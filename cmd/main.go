package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Masha003/Golang-gateway/internal/api"
	"github.com/Masha003/Golang-gateway/internal/api/controllers"
	"github.com/Masha003/Golang-gateway/internal/config"
	"github.com/Masha003/Golang-gateway/internal/service"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	// User
	userService, err := service.NewUserService(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	userController := controllers.NewUserController(userService)

	// Start HTTP Server
	httpSrv := api.NewHttpServer(cfg, userController)
	go func() {
		if err := httpSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Failed to start server: ", err)
		}
		log.Print("All server connections are closed")
	}()

	// Gracefull Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)
	<-quit
	log.Print("Shutting down server...")

	// Shutdown HTTP Server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	// Shurdown GRPC client
	err = userService.Close()
	if err != nil {
		log.Fatal("Failed to close GRPC client: ", err)
	}

	log.Print("Server exited properly")
}
