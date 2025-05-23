package main

import (
	c "breyting/blog-aggregator/internal/config"
	"breyting/blog-aggregator/internal/database"
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := c.Read()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		panic(err)
	}
	dbQueries := database.New(db)

	newState := &c.State{
		Db:     dbQueries,
		Config: &cfg,
	}

	commands := &c.Commands{
		Names: make(map[string]func(*c.State, c.Command) error),
	}

	commands.Register("login", c.HandlerLogin)
	commands.Register("register", c.HandlerRegister)
	commands.Register("reset", c.HandlerReset)
	commands.Register("users", c.HandlerGetUsers)
	commands.Register("agg", c.HandlerFetchFeed)
	commands.Register("addfeed", middlewareLoggedIn(c.HandlerAddFeed))
	commands.Register("feeds", c.HandlerGetFeeds)
	commands.Register("follow", middlewareLoggedIn(c.HandlerFollow))
	commands.Register("following", middlewareLoggedIn(c.HandlerFollowing))
	commands.Register("unfollow", middlewareLoggedIn(c.HandlerUnfollow))

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

func middlewareLoggedIn(handler func(s *c.State, cmd c.Command, user database.User) error) func(*c.State, c.Command) error {
	return func(s *c.State, cmd c.Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
