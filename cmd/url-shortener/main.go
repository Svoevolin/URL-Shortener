package main

import (
	"fmt"
	"url-shortener/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(*cfg)

	// TO DO: init logger: slog

	// TO DO: init storage: sqlite

	// TO DO: init router: chi, "chi render"

	// TO DO: run server
}
