package tcp

import (
	"context"
	"errors"
	"fmt"
	"github.com/addfs/go-mem-database/internal/config"
	"github.com/addfs/go-mem-database/internal/tools"
	"go.uber.org/zap"
	"net"
	"sync"
	"time"
)

const defaultMaxMessageSize = 2048
const defaultIdleTimeout = time.Minute * 5

type Handler = func(context.Context, []byte) []byte

type Server struct {
	address     string
	semaphore   tools.Semaphore
	idleTimeout time.Duration
	messageSize int
	logger      *zap.Logger
}

func NewServer(config *config.Config, logger *zap.Logger) (*Server, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	return &Server{
		address:     config.Network.Address,
		semaphore:   tools.NewSemaphore(config.Network.MaxConnections),
		idleTimeout: defaultIdleTimeout,
		messageSize: defaultMaxMessageSize,
		logger:      logger,
	}, nil
}

func (s *Server) HandleQueries(ctx context.Context, handler Handler) error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		for {
			conn, err := listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}

				s.logger.Error("failed to accept", zap.Error(err))
				continue
			}
			wg.Add(1)

			go func(c net.Conn) {
				s.semaphore.Acquire()

				defer func() {
					s.semaphore.Release()
					wg.Done()
				}()

				s.handleConnection(ctx, c, handler)

			}(conn)
		}
	}()

	go func() {
		defer wg.Done()

		<-ctx.Done()
		if err := listener.Close(); err != nil {
			s.logger.Warn("failed to close listener", zap.Error(err))
		}
	}()

	wg.Wait()

	return nil
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn, handler Handler) {
	defer func() {
		if err := conn.Close(); err != nil {
			s.logger.Error("failed to close connection", zap.Error(err))
		}
	}()

	request := make([]byte, s.messageSize)
	for {
		if err := conn.SetDeadline(time.Now().Add(s.idleTimeout)); err != nil {
			s.logger.Warn("failed to set read deadline", zap.Error(err))
			break
		}

		n, err := conn.Read(request)
		if err != nil {
			s.logger.Error("failed to read from connection", zap.Error(err))
			return
		}

		s.logger.Info("received message", zap.String("message", string(request[:n])))

		response := handler(ctx, request[:n])
		if _, err := conn.Write(response); err != nil {
			s.logger.Error("failed to write to connection", zap.Error(err))
			return
		}
	}
}
