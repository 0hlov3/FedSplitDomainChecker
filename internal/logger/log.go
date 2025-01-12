package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ZapLog *zap.Logger

func Init(verbose bool) {
	var err error
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false, // Set to true for development-style logs
		Sampling:    nil,
		Encoding:    "console", // Use "console" for human-readable logs
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder, // Add color to levels
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	// Adjust logging level for verbose mode
	if verbose {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	ZapLog, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}
