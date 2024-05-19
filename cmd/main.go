package main

import (
	"bufio"
	"context"
	"fmt"
	initialization "github.com/addfs/go-mem-database/internal"
	config2 "github.com/addfs/go-mem-database/internal/config"
	"github.com/addfs/go-mem-database/internal/database"
	"github.com/addfs/go-mem-database/internal/database/compute"
	"github.com/addfs/go-mem-database/internal/database/storage"
	"github.com/addfs/go-mem-database/internal/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Set up signal catching
	sigs := make(chan os.Signal, 1)

	// Catch all signals since we're not specifying which signal to catch
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Now the program will not exit and will wait for incoming signals
	// This is a goroutine which will keep on running in background
	go func() {
		s := <-sigs
		switch s {
		case syscall.SIGINT:
			fmt.Println("Received SIGINT, exiting...")
			os.Exit(0)
		case syscall.SIGTERM:
			fmt.Println("Received SIGTERM, exiting...")
			os.Exit(0)
		default:
			fmt.Println("Unknown signal.")
		}
	}()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yaml"
	}
	config := config2.NewConfig(configPath)()
	logger := log.NewLogger(config)

	parser, err := compute.NewParser(logger)
	if err != nil {

	}
	analyzer, err := compute.NewAnalyzer(logger)
	if err != nil {

	}
	computeLayer, err := compute.NewCompute(parser, analyzer, logger)
	if err != nil {

	}

	dbEngine, err := initialization.CreateEngine(config, logger)
	if err != nil {
		//return nil, fmt.Errorf("failed to initialize engine: %w", err)
	}

	storageLayer, err := storage.NewStorage(dbEngine, logger)
	if err != nil {
		//i.logger.Error("failed to initialize storage layer", zap.Error(err))
		//return nil, err
	}

	db, err := database.NewDatabase(computeLayer, storageLayer, logger)
	if err != nil {
		fmt.Println("Failed to initialize database:", err)
		os.Exit(1)
	}
	ctx := context.Background()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()
		result := db.HandleQuery(ctx, command)
		fmt.Println("Received result:", result)
	}
}
