package compose

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	Datacenters map[string]Datacenter `yaml:"datacenters"`
	Commands    Commands
}

type execResult struct {
	datacenterID string
	serviceID    string
	command      Command
}

// Exec ...
func (c *Compose) Exec(args []string) error {

	serviceCmdAlias := args[0]

	cmds, present := c.Commands[serviceCmdAlias]
	var err error
	if present {
		for _, cmd := range cmds {
			err = c.execServiceCmd(cmd)
		}
	} else {
		err = c.execServiceCmds(args)
	}
	return err
}

func (c *Compose) execServiceCmd(args string) error {
	return c.execServiceCmds(strings.Fields(args))
}

func (c *Compose) execServiceCmds(args []string) error {

	fmt.Println("Exec args:" + strings.Join(args, " "))

	// check if global command

	// else find datacenter service
	datacenterName := args[0]

	datacenter, present := c.Datacenters[datacenterName]
	if !present {
		return errors.New("Invalid datacenter name")
	}

	serviceName := args[1]
	service, present := datacenter.Services[serviceName]
	if !present {
		return errors.New("Invalid service name")
	}

	servicePath := filepath.Join(c.projectDir, service.Path)
	err := os.Chdir(servicePath)
	if err != nil {
		return err
	}

	command := args[2]
	commandArgs := args[3:]

	// search if a command is defined
	commandList, present := service.Commands[command]
	if present {
		commands := commandList[0]

		commandsSplit := strings.Fields(commands)
		return executeCommand(commandsSplit[0], commandsSplit[1:], servicePath, datacenter.Environment)
	}

	// Execute command in service directory
	return executeCommand(command, commandArgs, servicePath, datacenter.Environment)
}

func executeCommand(name string, args []string, dir string, env []string) error {
	cmd := exec.Command(name, args...)
	//	cmd.Args = args
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
