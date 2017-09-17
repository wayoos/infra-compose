package compose

import (
	"errors"
	"io/ioutil"
	"os"

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
	Version string
	//	projectDir  string
	Datacenters map[string]Datacenter `yaml:"datacenters"`
	Commands    Commands
}

// Exec ...
func (c *Compose) Exec(args []string) error {
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
