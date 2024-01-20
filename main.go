package main

import (
	"example.com/m/config"
	"example.com/m/internal/server"
	"log"
	"os"
)

func main() {
	cfg, err := config.MustConfig()

	if err != nil {
		cfg = &config.Config{Host: "127.0.0.1", Port: 8000}
	}

	runner := server.NewServer(cfg)
	if err := runner.Run(); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
