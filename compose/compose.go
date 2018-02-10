package compose

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"gopkg.in/src-d/go-git.v4"
	yaml "gopkg.in/yaml.v2"
)

type Command []string

type Commands map[string]Command

type VariableFile struct {
	File        string
	Environment Environment
}

// Service ...
type Service struct {
	Abstract    bool
	Parent      string
	Path        string
	Commands    Commands
	Command     Command
	Environment Environment
	Variables   map[string]VariableFile
}

type Environment []string

// Compose ... composed infrastructure
type Compose struct {
	Version    string
	projectDir string
	Services   map[string]Service
	//	Environments map[string]Environment
	Environment Environment

	DryRun bool
}

type execResult struct {
	//	environmentID string
	serviceID string
	commandID string
	command   Command
	execError error
}

type execResults struct {
	execResultList []*execResult
}

func (c *execResults) add(execResult execResult) {
	c.execResultList = append(c.execResultList, &execResult)
}

func newExecResults() *execResults {
	r := &execResults{}
	r.execResultList = make([]*execResult, 0)
	return r
}

// Exec ...
func (c *Compose) Exec(args []string) error {
	execResults := newExecResults()
	//serviceCmdAlias := args[0]

	// cmds, present := c.Commands[serviceCmdAlias]
	// var err error
	// if present {
	// 	for _, cmd := range cmds {
	// 		res := c.execServiceCmd(cmd)
	// 		execResults = append(execResults, res)
	// 		err = res.execError
	// 		if res.execError != nil {
	// 			break
	// 		}
	// 	}
	// } else {
	results, err := c.execServiceCmds(args, execResults)
	//execResults = append(execResults, res)
	//	err := res.execError
	//}

	dumpExecResults(results)

	return err
}

// List ... List all available command
func (c *Compose) List(args []string) error {
	const padding = 8
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "SERVICE\tCOMMAND\tSUB-COMMAND\t")

	var srvKeys []string
	for k, srv := range c.Services {
		if !srv.Abstract {
			srvKeys = append(srvKeys, k)
		}
	}
	sort.Strings(srvKeys)
	for _, srv := range srvKeys {
		service := c.Services[srv]

		if len(service.Command) > 0 {
			dumpCommandList(w, srv, "", service.Command)
		} else {
			commands := Commands{}

			for cmdKey, cmd := range service.Commands {
				commands[cmdKey] = cmd
			}

			dumpCommand(w, srv, commands)

		}

	}

	w.Flush()

	return nil
}

func (c *Compose) findServiceCommand(service Service) {

}

func dumpCommand(w *tabwriter.Writer, serviceName string, commands Commands) {
	// sort command
	var keys []string
	for k := range commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, cmd := range keys {
		subCommands := commands[cmd]

		dumpCommandList(w, serviceName, cmd, subCommands)
	}
}

func dumpCommandList(w *tabwriter.Writer, serviceName string, command string, subCommands []string) {
	commandList := ellipsis(40, strings.Join(subCommands, " | "))

	fmt.Fprintf(w, "%s\t%s\t%s\t\n", serviceName, command, commandList)
}

func ellipsis(length int, text string) string {
	r := []rune(text)
	if len(r) > length {
		return string(r[0:length]) + "..."
	}
	return text
}

func dumpExecResults(execResults *execResults) {
	fmt.Println("Execution summary")
	fmt.Println()
	const padding = 4
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "SERVICE\tCOMMAND\tSTATUS\t")
	for _, res := range execResults.execResultList {
		status := "Success"
		if res.execError != nil {
			status = "Error"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t\n",
			res.serviceID, res.commandID, status)
	}
	w.Flush()
	fmt.Println()
}

func (c *Compose) execServiceCmd(args string, execResults *execResults) (*execResults, error) {
	return c.execServiceCmds(strings.Fields(args), execResults)
}

func (c *Compose) execServiceCmds(args []string, execResults *execResults) (*execResults, error) {
	result := execResult{}
	//	fmt.Println("Exec args:" + strings.Join(args, " "))

	var env Environment

	// check if environment is defined
	//envID := args[0]
	//	envConf, present := c.Environments[envID]
	//	if present {
	//	env = envConf
	//args = args[1:]
	//result.environmentID = envID
	//	}

	serviceName := args[0]
	result.serviceID = serviceName
	service, present := c.Services[serviceName]
	if !present {
		return nil, errors.New("Invalid service name")
	}

	servicePath := filepath.Join(c.projectDir, service.Path)
	err := os.Chdir(servicePath)
	if err != nil {
		return nil, err
	}

	var command string
	var commandArgs []string
	if len(args) > 1 {
		command = args[1]
		if len(args) > 2 {
			commandArgs = args[2:]
		}
	}

	// Merge service environment
	env = appendEnv(service.Environment, env)

	// Find git branch name
	// TODO extract in util
	gitDir, gitErr := findGitRepo(servicePath)
	if gitErr == nil {
		repo, repoErr := git.PlainOpen(gitDir)
		if repoErr == nil {
			ref, headErr := repo.Head()
			if headErr == nil {
				branch := ref.Name().Short()

				branchSplit := strings.Split(branch, "/")

				branchFirst := branchSplit[0]
				branchLast := branchSplit[len(branchSplit)-1]

				os.Setenv("branch", branch)
				os.Setenv("branch.first", branchFirst)
				os.Setenv("branch.last", branchLast)

			}
		}
	}

	if !c.DryRun {
		// create variables files
		for _, variableFile := range service.Variables {
			fmt.Println("Var file  : " + variableFile.File)
			absProjectDir, _ := filepath.Abs(variableFile.File)
			parentDir := filepath.Dir(absProjectDir)

			fmt.Println("MkDir  : " + parentDir)
			os.MkdirAll(parentDir, 0755)

			outputVars := ""
			for _, variable := range variableFile.Environment {
				outputVars += os.ExpandEnv(variable) + "\n"
			}

			ioutil.WriteFile(variableFile.File, []byte(outputVars), 0644)
		}
	}

	if len(service.Command) > 0 {

		for _, commands := range service.Command {
			commandsSplit := strings.Fields(commands)

			cmd := commandsSplit[0]
			if strings.HasPrefix(cmd, "$") {
				newArgs := []string{cmd[1:]}
				args := commandsSplit[1:]
				newArgs = append(newArgs, args...)
				execResults, err := c.execServiceCmds(newArgs, execResults)
				if err != nil {
					return execResults, err
				}
			} else {
				result.commandID = cmd
				err = c.executeCommand(cmd, commandsSplit[1:], servicePath, env, service)
				if err != nil {
					result.execError = err
					execResults.add(result)
					return execResults, err
				}
				execResults.add(result)
			}

		}

		//		results := append(*execResults, result)
		return execResults, nil
	}

	// search if a command is defined
	commandList, present := service.Commands[command]
	if present {
		result.commandID = command
		for _, commands := range commandList {
			commandsSplit := strings.Fields(commands)
			cmd := commandsSplit[0]
			if strings.HasPrefix(cmd, "$") {
				newArgs := []string{cmd[1:]}
				args := commandsSplit[1:]
				newArgs = append(newArgs, args...)
				execResults, err := c.execServiceCmds(newArgs, execResults)
				if err != nil {
					return execResults, err
				}
			} else {
				err = c.executeCommand(cmd, commandsSplit[1:], servicePath, env, service)
				if err != nil {
					result.execError = err
					execResults.add(result)
					return execResults, err
				}
				execResults.add(result)
			}
		}
		return execResults, err
	}

	// Execute command in service directory
	result.commandID = "-"
	result.execError = c.executeCommand(command, commandArgs, servicePath, env, service)
	execResults.add(result)
	return execResults, err
}

func (c *Compose) executeCommand(name string, args []string, dir string, env Environment, service Service) error {
	os.Setenv("service.home", dir)

	argsExpandedEnv := []string{}
	for _, arg := range args {
		argsExpandedEnv = append(argsExpandedEnv, os.ExpandEnv(arg))
	}

	envExpandedEnv := []string{}
	for _, e := range env {
		envExpandedEnv = append(envExpandedEnv, os.ExpandEnv(e))
	}

	if c.DryRun {
		//		fmt.Println("Plan to Execute ")
		fmt.Println("Exec : " + name + " " + strings.Join(argsExpandedEnv, " "))
		fmt.Println("Dir  : " + dir)
		fmt.Println("Env  : " + strings.Join(envExpandedEnv, " "))
		fmt.Println("")
		return nil
	}

	cmd := exec.Command(name, argsExpandedEnv...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fullEnv := appendEnv(envExpandedEnv, os.Environ())
	cmd.Env = fullEnv

	err := cmd.Run()

	fmt.Println("State: " + cmd.ProcessState.String())

	if err != nil {
		return errors.New("Execute command error. " + err.Error())
	}

	return nil
}

// Load ...
func (c *Compose) Load(file string, projectDir string) error {
	validComposeFile, err := findComposeFile(file, projectDir)
	if err != nil {
		return err
	}

	err = c.loadCompose(validComposeFile)

	c.init()

	return err
}

func (c *Compose) mergeParent(service *Service, currentService Service) {
	if currentService.Parent != "" {
		parentService := c.Services[currentService.Parent]

		c.mergeParent(service, parentService)

		if parentService.Path != "" {
			service.Path = parentService.Path
		}

		for cmdKey, cmd := range parentService.Commands {
			if service.Commands == nil {
				service.Commands = make(Commands)
			}
			service.Commands[cmdKey] = cmd
		}

		for varKey, variable := range parentService.Variables {
			if service.Variables == nil {
				service.Variables = make(map[string]VariableFile)
			}
			service.Variables[varKey] = variable
		}

		for _, env := range parentService.Environment {
			if service.Environment == nil {
				service.Environment = Environment{}
			}
			service.Environment = append(service.Environment, env)
		}

	}

}

func (c *Compose) init() {
	services := make(map[string]Service)
	for serviceKey, service := range c.Services {
		if !service.Abstract {
			c.mergeParent(&service, service)
			for _, env := range c.Environment {
				if service.Environment == nil {
					service.Environment = Environment{}
				}
				service.Environment = append(service.Environment, env)
			}

			services[serviceKey] = service
		}
	}

	c.Services = services
}

func (c *Compose) loadCompose(composeFile string) error {
	source, err := ioutil.ReadFile(composeFile)
	if err != nil {
		return err
	}

	composeStr := string(source)

	//	os.Setenv("branch.first", "prod")
	//	os.Setenv("branch.last", "prod")

	//	composeParsed := os.ExpandEnv(composeStr)

	err = yaml.Unmarshal([]byte(composeStr), &c)
	if err != nil {
		return err
	}

	absFileName, _ := filepath.Abs(composeFile)
	c.projectDir = filepath.Dir(absFileName)

	return nil
}

func findGitRepo(dir string) (string, error) {
	gitDir := dir + "/.git"

	_, err := os.Stat(gitDir)
	if err != nil {
		// if not root path find in parent
		absProjectDir, _ := filepath.Abs(dir)
		parentDir := filepath.Dir(absProjectDir)

		if absProjectDir == "/" {
			return "", errors.New("Git repo not found")
		}

		return findGitRepo(parentDir)
	}

	return dir, nil
}

func findComposeFile(file string, projectDir string) (string, error) {
	if projectDir != "" {
		err := os.Chdir(projectDir)
		if err != nil {
			return "", err
		}
	}

	_, err := os.Stat(file)
	if err != nil {

		// if not root path find in parent
		absProjectDir, _ := filepath.Abs(projectDir)
		parentDir := filepath.Dir(absProjectDir)

		if absProjectDir == "/" {
			return "", errors.New("Compose file not found")
		}

		return findComposeFile(file, parentDir)
	}

	return file, nil
}
