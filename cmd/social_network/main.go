package main

import (
	"context"
	"os/signal"
	"social_network/internal/controllers"
	"social_network/internal/repositories"
	"social_network/internal/servers"
	"social_network/internal/services"
	"sync"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := initLogger()
	database := initDatabase("dsn", logger) // TODO: need to add dsn

	profileRepository := repositories.NewProfileRepository(database, logger)
	profileService := services.NewProfileService(&profileRepository)
	profileController := controllers.NewProfileController(&profileService, logger)

	restServer := servers.NewRESTServer("address", nil, &profileController, logger) // TODO: need to add address

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		restServer.ListenAndServe()
	}()

	go func() {
		defer wg.Done()
		<-ctx.Done()

		restServer.Shutdown(context.Background())
		logger.Info("REST server closed")

		_ = database.Close()
		logger.Info("database closed")
	}()

	wg.Wait()
}
