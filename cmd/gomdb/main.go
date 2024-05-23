package main

import (
	"context"
	"fmt"
	initialization "github.com/addfs/go-mem-database/internal"
	config2 "github.com/addfs/go-mem-database/internal/config"
	"github.com/addfs/go-mem-database/internal/database"
	"github.com/addfs/go-mem-database/internal/database/compute"
	"github.com/addfs/go-mem-database/internal/database/storage"
	"github.com/addfs/go-mem-database/internal/log"
	"github.com/addfs/go-mem-database/internal/net/tcp"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server, err := tcp.NewServer(config, logger)
	if err != nil {
		//return nil, fmt.Errorf("failed to initialize network: %w", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		server.HandleQueries(ctx, func(ctx context.Context, query []byte) []byte {
			response := db.HandleQuery(ctx, string(query))
			return []byte(response)
		})
	}()

	wg.Wait()
}
