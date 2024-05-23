package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/addfs/go-mem-database/internal/net/tcp"
	"github.com/addfs/go-mem-database/internal/tools"
	"go.uber.org/zap"
	"os"
	"syscall"
	"time"
)

func main() {
	address := flag.String("address", "localhost:3223", "Address of the gomdb")
	idleTimeout := flag.Duration("idle_timeout", time.Minute, "Idle timeout for connection")
	maxMessageSizeStr := flag.String("max_message_size", "4KB", "Max message size for connection")
	flag.Parse()

	logger, _ := zap.NewProduction()
	maxMessageSize, err := tools.ParseSize(*maxMessageSizeStr)
	if err != nil {
		logger.Fatal("failed to parse max message size", zap.Error(err))
	}

	client, err := tcp.NewClient(*address, maxMessageSize, *idleTimeout)
	if err != nil {
		logger.Fatal("failed to connect with server", zap.Error(err))
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()
		response, err := client.Send([]byte(command))
		if err != nil {
			if errors.Is(err, syscall.EPIPE) {
				logger.Fatal("connection was closed", zap.Error(err))
			}

			logger.Error("failed to send query", zap.Error(err))
		}
		fmt.Println("Received result:", string(response))
	}
}
