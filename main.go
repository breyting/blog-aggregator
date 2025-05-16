package main

import (
	c "breyting/blog-aggregator/internal/config"
)

func main() {
	c.Read()
	c.SetUser()
}
