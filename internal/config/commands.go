package config

import "fmt"

type State struct {
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
	if err := s.Config.SetUser(username); err != nil {
		return err
	}
	fmt.Printf("Current user set as %s\n", username)
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
