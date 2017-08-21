package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

// Service ...
type Service struct {
	Path     string
	Commands map[string]string
}

// Datacenter ...
type Datacenter struct {
	Environment []string
	Services    map[string]Service
}

// Config ... is a blabla
type Config struct {
	Version     string
	projectDir  string
	Datacenters map[string]Datacenter `yaml:"datacenters"`
}

func read() {
	fmt.Printf("Hello, world.\n")

	var config Config

	filename := "infra-compose.yml"
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Value: %#v\n", config)

	m := yaml.MapSlice{}
	err = yaml.Unmarshal(source, &m)
	if err != nil {
		panic(err)
	}

	//	fmt.Printf("Value: %#v\n", m)

	var serviceBastion Service
	serviceBastion.Path = "mgmt/services/bastion"

	var c Config
	c.Version = "34"

	var dsgra Datacenter

	dsgra.Services = make(map[string]Service)
	dsgra.Services["bastion"] = serviceBastion

	c.Datacenters = make(map[string]Datacenter)
	c.Datacenters["gra"] = dsgra

	out, err := yaml.Marshal(c)

	//	fmt.Printf("Value: %#v\n", out)
	ioutil.WriteFile("test.yml", out, 0644)

}

func findCompose(c *cli.Context) (Config, error) {

	if c.GlobalIsSet("project-directory") {
		projectDir := c.GlobalString("project-directory")
		err := os.Chdir(projectDir)
		if err != nil {
			var config Config
			return config, cli.NewExitError(err, 1)
		}
	}

	filename := c.GlobalString("file")

	_, err := os.Stat(filename)
	if err != nil {
		var config Config
		return config, cli.NewExitError("Compose file not found", 1)
	}

	return loadCompose(filename)
}

func loadCompose(filename string) (Config, error) {
	var config Config

	source, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
		return config, cli.NewExitError(err, 1)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		return config, cli.NewExitError(err, 1)
	}

	absFileName, _ := filepath.Abs(filename)
	config.projectDir = filepath.Dir(absFileName)

	return config, nil
}

func execCommand(c *cli.Context) error {

	config, err := findCompose(c)
	if err != nil {
		return err
	}

	if !c.Args().Present() {
		return cli.NewExitError("\"infra-compose exec\" requires at least one argument.", 1)
	}

	command := c.Args().First()

	// check if global command

	// else find datacenter service
	datacenter, present := config.Datacenters[command]
	if present {
		if c.NArg() >= 2 {

			var serviceArgs cli.Args
			serviceArgs = c.Args().Tail()

			service, present := datacenter.Services[serviceArgs.First()]
			if present {
				servicePath := filepath.Join(config.projectDir, service.Path)
				os.Chdir(servicePath)

				var commandArgs cli.Args
				commandArgs = serviceArgs.Tail()
				cmd := exec.Command(commandArgs.First())
				cmd.Args = commandArgs
				cmd.Dir = servicePath
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Stdin = os.Stdin

				err := cmd.Run()
				if err != nil {
					fmt.Printf("Failed to start. %s\n", err.Error())

					return cli.NewExitError("execute command error", 1)
				}

			} else {
				return cli.NewExitError("Invalid service", 1)
			}
		} else {
			return cli.NewExitError("Service name required", 1)
		}
	} else {
		return cli.NewExitError("Invalid datacenter", 1)
	}

	return nil
}

func main() {

	app := cli.NewApp()
	app.Name = "infra-compose"
	app.Usage = "Define and run infrastructure."
	app.UsageText = "infra-compose [global options] command [command options] [arguments...]"
	app.Version = "0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "file, f", Usage: "Specify an alternate compose file", Value: "infra-compose.yml"},
		cli.StringFlag{Name: "project-directory, p", Usage: "Specify an alternate working directory (default: the path of the Compose file)"},
	}

	app.Commands = []cli.Command{
		{
			Name:      "services",
			Aliases:   []string{"s"},
			Usage:     "List services",
			UsageText: "List all available services.",
			Action: func(c *cli.Context) error {
				fmt.Fprintf(c.App.Writer, "DATACENTER    SERVICES\n")
				fmt.Fprintf(c.App.Writer, "gra           global\n")
				fmt.Fprintf(c.App.Writer, "gra           bastion\n")
				fmt.Fprintf(c.App.Writer, "sbg           global\n")
				return nil
			},
		},
		{
			Name:      "commands",
			Aliases:   []string{"c"},
			Usage:     "List commands",
			UsageText: "List all available commands.",
			Action: func(c *cli.Context) error {
				fmt.Fprintf(c.App.Writer, "COMMANDS\n")
				fmt.Fprintf(c.App.Writer, "up\n")
				fmt.Fprintf(c.App.Writer, "down\n")
				return nil
			},
		},
		{
			Name:      "exec",
			Aliases:   []string{"e"},
			Usage:     "Run a global command or in a service",
			UsageText: "Run a global command or in a service",
			Action:    execCommand,
		},
	}

	app.Action = func(c *cli.Context) {

		if len(c.Args()) <= 0 {
			cli.ShowAppHelp(c)
		} else {
			println("Invalid command")
		}
	}

	app.Run(os.Args)

}
