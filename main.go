package main

import (
	c "breyting/blog-aggregator/internal/config"
	"fmt"
)

func main() {
	cfg, err := c.Read()
	if err != nil {
		panic(err)
	}

	err = cfg.SetUser("lane")
	if err != nil {
		panic(err)
	}

	cfg, err = c.Read()
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg)
}
