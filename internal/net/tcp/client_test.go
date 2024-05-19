package tcp

import (
	"github.com/stretchr/testify/require"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestClient_Send(t *testing.T) {
	t.Parallel()

	request := []byte("hello server")
	response := []byte("hello client")

	listener, err := net.Listen("tcp", ":10001")
	require.NoError(t, err)

	go func() {
		connection, err := listener.Accept()
		if err != nil {
			return
		}

		buffer := make([]byte, 2048)
		count, err := connection.Read(buffer)
		require.NoError(t, err)
		require.Equal(t, request, buffer[:count])

		_, err = connection.Write(response)
		require.NoError(t, err)

		defer func() {
			err = connection.Close()
			require.NoError(t, err)
			err = listener.Close()
			require.NoError(t, err)
		}()
	}()

	client, err := NewClient("127.0.0.1:10001", 2048, time.Minute)
	require.NoError(t, err)

	buffer, err := client.Send(request)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual([]byte(response), buffer))
}
