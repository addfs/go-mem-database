package initialization

import (
	"fmt"
	config2 "github.com/addfs/go-mem-database/internal/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"path"
	"runtime"
	"testing"
)

func TestCreateEngineWithoutConfig(t *testing.T) {
	t.Parallel()

	engine, err := CreateEngine(nil, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, engine)
}

func TestCreateEngine(t *testing.T) {
	t.Parallel()

	_, file, _, _ := runtime.Caller(0)
	rootPath := path.Clean(file + "/../../")

	t.Run("TestCreateEngine", func(t *testing.T) {
		config := config2.NewConfig(fmt.Sprintf("%s/config/config.yaml", rootPath))()
		engine, err := CreateEngine(config, zap.NewNop())
		require.NoError(t, err)
		require.NotNil(t, engine)
	})

	t.Run("TestCreateEngine", func(t *testing.T) {
		config := config2.NewConfig(fmt.Sprintf("%s/config/config.yaml", rootPath))()
		config.EngineConfig.Type = "incorrect"
		engine, err := CreateEngine(config, zap.NewNop())
		require.Error(t, err, "engine type is incorrect")
		require.Nil(t, engine)
	})

}
