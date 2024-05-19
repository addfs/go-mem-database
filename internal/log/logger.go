package log

import (
	"github.com/addfs/go-mem-database/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func NewLogger(config *config.Config) *zap.Logger {
	logger, err := zap.Config{
		Level:       config.Logger.Level,
		Development: false,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "message",
			LevelKey:      "level",
			TimeKey:       "timestamp",
			NameKey:       "logger",
			CallerKey:     "caller",
			FunctionKey:   zapcore.OmitKey,
			StacktraceKey: "",
			LineEnding:    "\n",
			EncodeLevel: func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
				encoder.AppendString("level=" + level.String())
			},
			EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
				encoder.AppendString("time=" + t.Format(time.RFC3339))
			},
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller: func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
				encoder.AppendString("caller=" + caller.TrimmedPath())
			},
			EncodeName:       zapcore.FullNameEncoder,
			ConsoleSeparator: " ",
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()

	if err != nil {
		panic(err)
	}

	return logger
}
