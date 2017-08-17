package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

type Service struct {
	Path     string
	Commands map[string]string
}

type Datacenter struct {
	Environment []string
	Services    map[string]Service
}

// Config ... is a blabla
type Config struct {
	Version     string
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

func main() {

	app := cli.NewApp()
	app.Name = "infra-compose"
	app.Usage = "Define and run infrastructure."
	app.UsageText = "infra-compose [global options] command [command options] [arguments...]"
	app.Version = "0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "file, f", Usage: "Specify an alternate compose file", Value: "infra-compose.yml"},
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
			Action: func(c *cli.Context) error {

				fmt.Printf("Hello %q", c.Args().Get(0))
				fmt.Fprintf(c.App.Writer, "COMMANDS\n")
				fmt.Fprintf(c.App.Writer, "up\n")
				fmt.Fprintf(c.App.Writer, "down\n")
				return nil
			},
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
