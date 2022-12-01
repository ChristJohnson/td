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

	args := c.Args().Slice()

	for arg := range args {
		if args[arg] == "init" {
			return nil
		}
	}

	if path == "" {
		fmt.Println()
		ct.ChangeColor(ct.Red, false, ct.None, false)
		fmt.Println("Error")
		fmt.Println("-----")
		ct.ResetColor()

		errormessage := "A file for your todos is missing!\n" +
			"If you see this message, then the program has attempted to create\n" +
			"a file in the current directory named .todos. If this is not what\n" +
			"you intended, then set the environment variable, TODO_DB_PATH, to\n" +
			"the directory you would like to hold the file."

		fmt.Println(errormessage)
	}

	path, err = os.Getwd()
	if err != nil {
		fmt.Println("There was an error making the .todos file!")
	} else {
		err = createStoreFileIfNeeded(path + "/.todos")
	}

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
