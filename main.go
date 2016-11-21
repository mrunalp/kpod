package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

const (
	version = "0.0.1"
)

func main() {
	app := cli.NewApp()
	app.Name = "kpod"
	app.Usage = "Manage pods and images"
	app.Version = version
	app.Author = "Kpod Contributors"
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "path",
			Value: "",
			Usage: "Path to runc bundle to run",
		},
	}
	app.Commands = []cli.Command{}

	app.Action = func(c *cli.Context) error {
		if c.String("path") == "" {
			log.Fatal("Must specify path to runc bundle to run")
		}

		// TODO use a reasonable containerID - randomly generate or pull from a flag
		if err := RunContainer("123456", c.String("path")); err != nil {
			return err
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
