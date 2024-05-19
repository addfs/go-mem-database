package initialization

import (
	"errors"
	"github.com/addfs/go-mem-database/internal/config"
	"github.com/addfs/go-mem-database/internal/database/storage"
	"github.com/addfs/go-mem-database/internal/database/storage/engine/in_memory"
	"go.uber.org/zap"
)

const (
	InMemoryEngine = "in_memory"
)

var supportedEngineTypes = map[string]struct{}{
	InMemoryEngine: {},
}

const defaultPartitionsNumber = 10

func CreateEngine(cfg *config.Config, logger *zap.Logger) (storage.Engine, error) {
	if cfg == nil {
		return in_memory.NewEngine(in_memory.HashTableBuilder, defaultPartitionsNumber, logger)
	}

	engineConfig := cfg.EngineConfig

	if engineConfig.Type != "" {
		_, found := supportedEngineTypes[engineConfig.Type]
		if !found {
			return nil, errors.New("engine type is incorrect")
		}
	}

	return in_memory.NewEngine(in_memory.HashTableBuilder, defaultPartitionsNumber, logger)
}
