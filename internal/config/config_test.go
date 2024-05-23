package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"path"
	"runtime"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()

	_, file, _, _ := runtime.Caller(0)
	rootPath := path.Clean(file + "/../../../")

	t.Run("file", func(t *testing.T) {
		c := NewConfig(fmt.Sprintf("%s/test/testdata/config.yaml", rootPath))()

		// Logger
		assert.Equal(t, "info", c.Logger.Level.String(), "logger.level")
		assert.Equal(t, "/log/output.log", c.Logger.Output, "logger.output")

		// Engine
		assert.Equal(t, "in_memory", c.EngineConfig.Type, "engine.type")

		// Network
		assert.Equal(t, "127.0.0.1:4444", c.Network.Address, "network.address")
		assert.Equal(t, 50, c.Network.MaxConnections, "network.max_connections")
	})

	t.Run("default", func(t *testing.T) {
		c := NewConfig(fmt.Sprintf("%s/test/testdata/empty_config.yaml", rootPath))()

		// Logger
		assert.Equal(t, "info", c.Logger.Level.String(), "logger.level")
		assert.Equal(t, "", c.Logger.Output, "logger.output")

		// Engine
		assert.Equal(t, "in_memory", c.EngineConfig.Type, "engine.type")

		// Network
		assert.Equal(t, "127.0.0.1:3223", c.Network.Address, "network.address")
		assert.Equal(t, 100, c.Network.MaxConnections, "network.max_connections")
	})
}
