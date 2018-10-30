package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tonto/kit/goapp/pkg/command"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "Goapp CLI"
	app.Version = "0.1.0"

	app.Commands = []cli.Command{
		cli.Command{
			Name:      "new",
			Usage:     "new [--app | --svc] name",
			UsageText: "creates app or service",
			Description: `Create new application or service
			--app -a 	Creates new application (default behavior) using provided application name 
			--svc -s	Creates new service for existing application. Must be inside application directory 
	  `,
			SkipFlagParsing: false,
			HideHelp:        false,
			Hidden:          false,
			ArgsUsage:       "[-a projname] [-s svcname]",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "app, a"},
				cli.BoolFlag{Name: "svc, s"},
				// TODO - other options that define app / service templates eg. -tlm (metrics, logging, transport) -t http
			},
			HelpName: "new",

			// TODO - Add coloring

			Before: func(c *cli.Context) error {
				if len(c.Args()) == 0 {
					return fmt.Errorf("you must provide app/service name")
				}

				return nil
			},
			Action: newCmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func newCmd(c *cli.Context) error {
	var cmd command.Command

	cmd, err := command.NewCreateApp(c.Args()[0], "default.json")
	if err != nil {
		return err
	}

	return execCmd(cmd)

	// c.Command.FullName()
	// c.Command.VisibleFlags()
	// fmt.Fprintf(c.App.Writer, "here \n")
	// if c.Bool("svc") {
	// 	fmt.Fprintf(c.App.Writer, "creating service\n")
	// 	return nil
	// }

	// fmt.Fprintf(c.App.Writer, "creating app\n")

	// return nil
}

func execCmd(cmd command.Command) (err error) {
	defer func() {
		if err != nil {
			fmt.Printf("%v - rolling back.", err)
			err = cmd.Rollback()
		}
	}()
	return cmd.Execute()
}
