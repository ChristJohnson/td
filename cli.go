package main

import (
	"fmt"
	"os"

	ct "github.com/daviddengcn/go-colortext"
	cli "github.com/urfave/cli/v2"
)

type Error string

func (e Error) Error() string { return string(e) }

const argError = Error("Error in argument")

func main() {

	/* TODO: finish writing out the declarative approach:
	https://cli.urfave.org/v2/examples/full-api-example/#
	https://github.com/Swatto/td
	https://github.com/urfave/cli/tree/v2-maint
	*/

	/* For the tac behavior of standard cli tools

	ex:
	td -a
	td --all
	td -d --all
	*/
	flags := []cli.Flag{
		&cli.BoolFlag{
			Name:  "done, d",
			Usage: "print done todos",
		},
		&cli.BoolFlag{
			Name:  "all, a",
			Usage: "print all todos",
		},
	}

	app := &cli.App{
		Name:      "td",
		HelpName:  "td",
		Usage:     "todo list manager",
		UsageText: "td - manage your todo list",
		Version:   "1.4.2",
		Flags:     flags,
		Before:    action_before,
		Action:    action_default,
		Commands:  cli_commands_list(),
		Authors: []*cli.Author{
			{
				Name:  "Gaël Gillard",
				Email: "gillardgael@gmail.com",
			},
			{
				Name: "ChristJohnson",
			},
		},
	}

	app.Run(os.Args)
}

func action_before(c *cli.Context) error {
	var err error
	path := getDBPath()

	if path == "" {
		fmt.Println()
		ct.ChangeColor(ct.Red, false, ct.None, false)
		fmt.Println("Error")
		fmt.Println("-----")
		ct.ResetColor()
		fmt.Println("A store for your todos is missing. You have 2 possibilities:")
		fmt.Println("  1. create a \".todos\" file in your local folder.")
		fmt.Println("  2. the environment variable \"TODO_DB_PATH\" could be set.")
		fmt.Println("    (example: \"export TODO_DB_PATH=$HOME/Dropbox/todo.json\" in your .bashrc or .bash_profile)")
		fmt.Println()
	}

	createStoreFileIfNeeded(path)

	return err
}

func action_default(c *cli.Context) error {
	var err error
	collection := collection{}

	err = collection.RetrieveTodos()
	if err != nil {
		fmt.Println(err)
	} else {
		if !c.IsSet("all") {
			if c.IsSet("done") {
				collection.ListDoneTodos()
			} else {
				collection.ListPendingTodos()
			}
		}

		if len(collection.Todos) > 0 {
			fmt.Println()
			for _, todo := range collection.Todos {
				todo.MakeOutput(true)
			}
			fmt.Println()
		} else {
			ct.ChangeColor(ct.Cyan, false, ct.None, false)
			fmt.Println("There's no todo to show.")
			ct.ResetColor()
		}
	}
	return nil
}
