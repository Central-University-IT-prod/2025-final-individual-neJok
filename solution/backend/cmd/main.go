package main

import (
	"errors"
	"fmt"
	"log"
	"neJok/solution/app"
	"neJok/solution/config"
	_ "neJok/solution/pkg/docs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// @title PROD Backend 2025 Advertising Platform API
// @version 1.0
func main() {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	fmt.Printf("Server started")
	go func() {
		err := router.Run(cfg.ServerAddress)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("run server: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sig := <-quit

	fmt.Printf("Received signal: %s\n", sig)
}
