package main

import (
	"fmt"
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

	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{}

	app.Run(os.Args)
}

func generateBashCompletion(c *cli.Context) {
	tasks := []string{}

	// This will complete if no args are passed
	if c.NArg() > 0 {
		return
	}
	for _, t := range tasks {
		fmt.Println(t)
	}
}
