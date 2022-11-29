package main

import (
	"fmt"
	"os"
	"strconv"

	ct "github.com/daviddengcn/go-colortext"
	cli "github.com/urfave/cli/v2"
)

func cli_commands_list() cli.Commands {
	return cli.Commands{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Initialize a collection of todos",
			Action:  cli_init,
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add a new todo",
			Action:  cli_add,
		},
		{
			Name:    "modify",
			Aliases: []string{"m"},
			Usage:   "Modify the text of an existing todo",
			Action:  cli_modify,
		},
		{
			Name:    "toggle",
			Aliases: []string{"t"},
			Usage:   "Toggle the status of a todo by giving his id",
			Action:  cli_toggle,
		},
		{
			Name:    "clean",
			Aliases: []string{"c"},
			Usage:   "Remove finished todos from the list",
			Action:  cli_clean,
		},
		{
			Name:    "reorder",
			Aliases: []string{"r"},
			Usage:   "Reset ids of todo (no arguments) or swap the position of two todos",
			Action:  cli_reorder,
		},
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "Search a string in all todos",
			Action:  cli_search,
		},
	}
}

func cli_init(*cli.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("%s .\n", err)
		return err
	}

	err = createStoreFileIfNeeded(cwd + "/.todos")
	ct.ChangeColor(ct.Cyan, false, ct.None, false)
	if err != nil {
		fmt.Printf("A \".todos\" file already exist in \"%s\".\n", cwd)
	} else {
		fmt.Printf("A \".todos\" file is now added to \"%s\".\n", cwd)
	}
	ct.ResetColor()
	return nil
}

func cli_add(c *cli.Context) error {

	if c.Args().Len() != 1 {
		fmt.Println()
		ct.ChangeColor(ct.Red, false, ct.None, false)
		fmt.Println("Error")
		ct.ResetColor()
		fmt.Println("You must provide a name to your todo.")
		fmt.Println("Example: td add \"call mum\"")
		fmt.Println()
		return argError
	}

	collection := collection{}
	todo := todo{
		ID:       0,
		Desc:     c.Args().Get(0),
		Status:   "pending",
		Modified: "",
	}
	err := collection.RetrieveTodos()
	if err != nil {
		fmt.Println(err)
		return err
	}

	id, err := collection.CreateTodo(&todo)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ct.ChangeColor(ct.Cyan, false, ct.None, false)
	fmt.Printf("#%d \"%s\" is now added to your todos.\n", id, c.Args().Get(0))
	ct.ResetColor()
	return nil
}

func cli_modify(c *cli.Context) error {

	if c.Args().Len() != 2 {
		fmt.Println()
		ct.ChangeColor(ct.Red, false, ct.None, false)
		fmt.Println("Error")
		ct.ResetColor()
		fmt.Println("You must provide the id and the new text for your todo.")
		fmt.Println("Example: td modify 2 \"call dad\"")
		fmt.Println()
		return argError
	}

	collection := collection{}
	collection.RetrieveTodos()

	args := c.Args()

	id, err := strconv.ParseInt(args.Get(0), 10, 32)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = collection.Modify(id, args.Get(1))
	if err != nil {
		fmt.Println(err)
		return err
	}

	ct.ChangeColor(ct.Cyan, false, ct.None, false)
	fmt.Printf("\"%s\" has now a new description: %s\n", args.Get(0), args.Get(1))
	ct.ResetColor()
	return nil
}

func cli_toggle(c *cli.Context) error {
	var err error

	if c.Args().Len() != 1 {
		fmt.Println()
		ct.ChangeColor(ct.Red, false, ct.None, false)
		fmt.Println("Error")
		ct.ResetColor()
		fmt.Println("You must provide the position of the item you want to change.")
		fmt.Println("Example: td toggle 1")
		fmt.Println()
		return argError
	}

	collection := collection{}
	collection.RetrieveTodos()

	id, err := strconv.ParseInt(c.Args().Get(0), 10, 32)
	if err != nil {
		fmt.Println(err)
		return err
	}

	todo, err := collection.Toggle(id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ct.ChangeColor(ct.Cyan, false, ct.None, false)
	fmt.Printf("Your todo is now %s.\n", todo.Status)
	ct.ResetColor()
	return nil
}

func cli_clean(c *cli.Context) error {
	collection := collection{}
	collection.RetrieveTodos()

	err := collection.RemoveFinishedTodos()

	if err != nil {
		fmt.Println(err)
		return err
	}

	ct.ChangeColor(ct.Cyan, false, ct.None, false)
	fmt.Println("Your list is now flushed of finished todos.")
	ct.ResetColor()
	return nil
}

func cli_reorder(c *cli.Context) error {
	collection := collection{}
	collection.RetrieveTodos()

	if c.Args().Len() != 1 {
		fmt.Println()
		ct.ChangeColor(ct.Red, false, ct.None, false)
		fmt.Println("Error")
		ct.ResetColor()
		fmt.Println("You must provide two position if you want to swap todos.")
		fmt.Println("Example: td reorder 9 3")
		fmt.Println()
		return argError
	} else if c.Args().Len() != 2 {
		idA, err := strconv.ParseInt(c.Args().Get(0), 10, 32)
		if err != nil {
			fmt.Println(err)
			return err
		}

		idB, err := strconv.ParseInt(c.Args().Get(1), 10, 32)
		if err != nil {
			fmt.Println(err)
			return err
		}

		_, err = collection.Find(idA)
		if err != nil {
			fmt.Println(err)
			return err
		}

		_, err = collection.Find(idB)
		if err != nil {
			fmt.Println(err)
			return err
		}

		collection.Swap(idA, idB)

		ct.ChangeColor(ct.Cyan, false, ct.None, false)
		fmt.Printf("\"%s\" and \"%s\" has been swapped\n", c.Args().Get(0), c.Args().Get(1))
		ct.ResetColor()
	}

	err := collection.Reorder()

	if err != nil {
		fmt.Println(err)
		return err
	}

	ct.ChangeColor(ct.Cyan, false, ct.None, false)
	fmt.Println("Your list is now reordered.")
	ct.ResetColor()
	return nil
}

func cli_search(c *cli.Context) error {
	if c.Args().Len() != 1 {
		fmt.Println()
		ct.ChangeColor(ct.Red, false, ct.None, false)
		fmt.Println("Error")
		ct.ResetColor()
		fmt.Println("You must provide a string earch.")
		fmt.Println("Example: td search \"project-1\"")
		fmt.Println()
		return argError
	}

	collection := collection{}
	collection.RetrieveTodos()
	collection.Search(c.Args().Get(0))

	if len(collection.Todos) == 0 {
		ct.ChangeColor(ct.Cyan, false, ct.None, false)
		fmt.Printf("Sorry, there's no todos containing \"%s\".\n", c.Args().Get(0))
		ct.ResetColor()
		return argError
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
	return nil
}
