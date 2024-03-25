package main

import (
	"fmt"

	"github.com/JamshedJ/URL-Shortener/internal/config"
)

func main() {
	// TODO: init config
	cfg := config.MustLoad()
	fmt.Println(cfg)
	// TODO: init logger
	// TODO: init storage
	// TODO: init router
	// TODO: run server
}
