package config

import (
	"breyting/blog-aggregator/internal/database"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

type State struct {
	Db     *database.Queries
	Config *Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Names map[string]func(*State, Command) error
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("login command requires exactly one argument")
	}
	username := cmd.Args[0]
	_, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		os.Exit(1)
	}

	if err := s.Config.SetUser(username); err != nil {
		return err
	}
	fmt.Printf("Current user set as %s\n", username)
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("login command requires exactly one argument")
	}

	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	})
	if err != nil {
		os.Exit(1)
	}

	if err := s.Config.SetUser(user.Name); err != nil {
		return err
	}
	fmt.Printf("New user created and current user set as %s\n", user.Name)
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	err := s.Db.Reset(context.Background())
	if err != nil {
		os.Exit(1)
	}
	return nil

}

func HandlerGetUsers(s *State, cmd Command) error {
	res, err := s.Db.GetUsers(context.Background())
	if err != nil {
		os.Exit(1)
	}

	for _, val := range res {
		if val.Name == s.Config.CurrentUserName {
			fmt.Printf("- %s (current)\n", val.Name)
		} else {
			fmt.Printf("- %s\n", val.Name)
		}
	}
	return nil
}

func (c *Commands) Run(s *State, cmd Command) error {
	if handler, ok := c.Names[cmd.Name]; ok {
		return handler(s, cmd)
	}
	return fmt.Errorf("unknown command: %s", cmd.Name)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	if c.Names == nil {
		c.Names = make(map[string]func(*State, Command) error)
	}
	c.Names[name] = f
}
