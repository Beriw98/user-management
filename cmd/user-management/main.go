package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/gommon/log"

	"github.com/Beriw98/user-management/internal/config"
	"github.com/Beriw98/user-management/internal/container"
	"github.com/Beriw98/user-management/internal/infrastructure/httpsrv"
)

func main() {
	cfg := config.New()

	ctr, err := container.NewContainer(cfg)
	if err != nil {
		panic(err)
	}

	r := httpsrv.NewRouter(ctr)

	go func() {
		if err = r.Start(":8080"); err != nil {
			log.Info("shutting down the server, message: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
	log.Info("initiating graceful shutdown...")

	ctx := context.Background()

	if err = r.Shutdown(ctx); err != nil {
		r.Logger.Fatal(err)
	}

}
