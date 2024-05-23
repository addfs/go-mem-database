package tcp

import (
	"context"
	"fmt"
	"github.com/addfs/go-mem-database/internal/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net"
	"path"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func TestTCPServer(t *testing.T) {
	t.Parallel()

	request := "hello server"
	response := "hello client"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := getConfig()
	server, err := NewServer(cfg, zap.NewNop())
	require.NoError(t, err)

	go func() {
		require.NoError(t, server.HandleQueries(ctx, func(ctx context.Context, buffer []byte) []byte {
			require.Equal(t, request, string(buffer))
			return []byte(response)
		}))
	}()
	time.Sleep(100 * time.Millisecond)

	connection, err := net.Dial("tcp", cfg.Network.Address)
	require.NoError(t, err)

	_, err = connection.Write([]byte(request))
	require.NoError(t, err)
	buffer := make([]byte, 2048)
	count, err := connection.Read(buffer)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual([]byte(response), buffer[:count]))
}

func getConfig() *config.Config {
	_, file, _, _ := runtime.Caller(0)
	rootPath := path.Clean(file + "/../../../../")
	return config.NewConfig(fmt.Sprintf("%s/config/config_test.yaml", rootPath))()
}
