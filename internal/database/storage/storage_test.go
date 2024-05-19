package storage

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"testing"
)

// mockgen -source=storage.go -destination=storage_mock.go -package=storage

func TestNewStorage(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	engine := NewMockEngine(ctrl)

	storage, err := NewStorage(nil, nil)
	require.Error(t, err, "engine is invalid")
	require.Nil(t, storage)

	storage, err = NewStorage(engine, nil)
	require.Error(t, err, "logger is invalid")
	require.Nil(t, storage)

	storage, err = NewStorage(engine, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, storage)
}

func TestSuccessfulSet(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "tx", int64(555))

	result := make(chan error, 1)
	result <- nil

	ctrl := gomock.NewController(t)
	engine := NewMockEngine(ctrl)
	engine.EXPECT().
		Set(ctx, "key", "value")

	storage, err := NewStorage(engine, zap.NewNop())
	require.NoError(t, err)

	err = storage.Set(ctx, "key", "value")
	require.NoError(t, err)
}

func TestSuccessfulGet(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "tx", int64(555))

	result := make(chan error, 1)
	result <- nil

	ctrl := gomock.NewController(t)
	engine := NewMockEngine(ctrl)
	engine.EXPECT().
		Get(ctx, "key").Return("value", true)

	storage, err := NewStorage(engine, zap.NewNop())
	require.NoError(t, err)

	value, err := storage.Get(ctx, "key")
	require.NoError(t, err)
	require.Equal(t, "value", value)
}

func TestSuccessfulDel(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "tx", int64(555))

	result := make(chan error, 1)
	result <- nil

	ctrl := gomock.NewController(t)
	engine := NewMockEngine(ctrl)
	engine.EXPECT().
		Del(ctx, "key")

	storage, err := NewStorage(engine, zap.NewNop())
	require.NoError(t, err)

	err = storage.Del(ctx, "key")
	require.NoError(t, err)
}
