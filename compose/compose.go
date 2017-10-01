package compose

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/tabwriter"

	yaml "gopkg.in/yaml.v2"
)

type Command []string

type Commands map[string]Command

// Service ...
type Service struct {
	Path     string
	Commands Commands
}

// Datacenter ...
type Datacenter struct {
	Environment []string
	Services    map[string]Service
}

// Compose ... composed infrastructure
type Compose struct {
	Version     string
	projectDir  string
	DryRun      bool
	Datacenters map[string]Datacenter `yaml:"datacenters"`
	Commands    Commands
}

type execResult struct {
	datacenterID string
	serviceID    string
	commandID    string
	command      Command
	execError    error
}

// Exec ...
func (c *Compose) Exec(args []string) error {
	var execResults []execResult
	serviceCmdAlias := args[0]

	cmds, present := c.Commands[serviceCmdAlias]
	var err error
	if present {
		for _, cmd := range cmds {
			res := c.execServiceCmd(cmd)
			execResults = append(execResults, res)
			err = res.execError
			if res.execError != nil {
				break
			}
		}
	} else {
		res := c.execServiceCmds(args)
		execResults = append(execResults, res)
		err = res.execError
	}

	dumpExecResults(execResults)

	return err
}

func dumpExecResults(execResults []execResult) {
	fmt.Println("Execution summary")
	fmt.Println()
	const padding = 4
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "Datacenter\tService\tCommand\tStatus\t")
	fmt.Fprintln(w, "\t\t\t\t\t")
	for _, res := range execResults {
		status := "Success"
		if res.execError != nil {
			status = "Error"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", res.datacenterID,
			res.serviceID, res.commandID, status)
	}
	w.Flush()
	fmt.Println()
}

func (c *Compose) execServiceCmd(args string) execResult {
	return c.execServiceCmds(strings.Fields(args))
}

func (c *Compose) execServiceCmds(args []string) execResult {
	result := execResult{}
	//	fmt.Println("Exec args:" + strings.Join(args, " "))

	// else find datacenter service
	datacenterName := args[0]

	result.datacenterID = datacenterName

	datacenter, present := c.Datacenters[datacenterName]
	if !present {
		result.execError = errors.New("Invalid datacenter name")
		return result
	}

	serviceName := args[1]
	result.serviceID = serviceName
	service, present := datacenter.Services[serviceName]
	if !present {
		result.execError = errors.New("Invalid service name")
		return result
	}

	servicePath := filepath.Join(c.projectDir, service.Path)
	err := os.Chdir(servicePath)
	if err != nil {
		result.execError = err
		return result
	}

	command := args[2]
	commandArgs := args[3:]

	// search if a command is defined
	commandList, present := service.Commands[command]
	if present {
		result.commandID = command
		for _, commands := range commandList {
			commandsSplit := strings.Fields(commands)
			c.executeCommand(commandsSplit[0], commandsSplit[1:], servicePath, datacenter.Environment)
		}
		return result
	}

	// Execute command in service directory
	result.commandID = "-"
	result.execError = c.executeCommand(command, commandArgs, servicePath, datacenter.Environment)
	return result
}

func (c *Compose) executeCommand(name string, args []string, dir string, env []string) error {
	if c.DryRun {
		fmt.Println("Dry-run: Plan to Execute ")
		fmt.Println("Exec : " + name + " " + strings.Join(args, " "))
		fmt.Println("Dir  : " + dir)
		fmt.Println("Env  : " + strings.Join(env, " "))
		fmt.Println("")
		return nil
	}

	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fullEnv := append(os.Environ(), env...)
	cmd.Env = fullEnv

	err := cmd.Run()
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

	return c.loadCompose(validComposeFile)
}

func (c *Compose) loadCompose(composeFile string) error {
	source, err := ioutil.ReadFile(composeFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(source, &c)
	if err != nil {
		return err
	}

	absFileName, _ := filepath.Abs(composeFile)
	c.projectDir = filepath.Dir(absFileName)

	return nil
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
