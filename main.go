package main

import (
	c "breyting/blog-aggregator/internal/config"
	"fmt"
	"os"
)

func main() {
	cfg, err := c.Read()
	if err != nil {
		panic(err)
	}

	newState := &c.State{
		Config: &cfg,
	}

	commands := &c.Commands{
		Names: make(map[string]func(*c.State, c.Command) error),
	}

	commands.Register("login", c.HandlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("Error: command must have at least 2 arguments")
		os.Exit(1)
	}
	cmd := c.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	err = commands.Run(newState, cmd)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
