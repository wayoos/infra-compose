package compose

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli"
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

// Exec ...
func (c *Compose) Exec(args []string) error {
	//fmt.Println("Exec args:" + strings.Join(args, " "))

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
	os.Chdir(servicePath)

	var commandArgs cli.Args
	commandArgs = args[2:]
	cmd := exec.Command(commandArgs.First())
	cmd.Args = commandArgs
	cmd.Dir = servicePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	env := append(os.Environ(), datacenter.Environment...)
	cmd.Env = env

	err := cmd.Run()
	if err != nil {
		return cli.NewExitError("Execute command error: "+err.Error(), 1)
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
		return "", errors.New("Compose file not found")
	}

	return file, nil
}
