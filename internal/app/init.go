package app

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Init zap logger as main logger.
func InitLog() zap.SugaredLogger {
	cfg := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stdout", "/tmp/gkc2.log"},
		ErrorOutputPaths: []string{"stdout", "/tmp/gkc2.log"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.RFC3339TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	zapLogger := zap.Must(cfg.Build())
	zapLogger.Info("logger construction succeeded")
	zap.ReplaceGlobals(zapLogger)
	defer func() {
		_ = zapLogger.Sync()
	}()

	sugar := *zapLogger.Sugar()

	defer func() {
		_ = sugar.Sync()
	}()
	return sugar
}
