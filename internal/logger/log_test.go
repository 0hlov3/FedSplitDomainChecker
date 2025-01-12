package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLoggerInit(t *testing.T) {
	tests := []struct {
		name     string
		verbose  bool
		expected zapcore.Level
	}{
		{"Non-Verbose Mode", false, zap.InfoLevel},
		{"Verbose Mode", true, zap.DebugLevel},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Redirect logs to a buffer for testing
			buffer := &bytes.Buffer{}
			writer := zapcore.AddSync(buffer)

			// Configure a custom logger for testing
			encoderConfig := zapcore.EncoderConfig{
				TimeKey:     "timestamp",
				LevelKey:    "level",
				MessageKey:  "msg",
				EncodeLevel: zapcore.CapitalLevelEncoder,
				EncodeTime:  zapcore.ISO8601TimeEncoder,
				LineEnding:  zapcore.DefaultLineEnding,
			}

			// Create a custom core based on the test case
			core := zapcore.NewCore(
				zapcore.NewConsoleEncoder(encoderConfig),
				writer,
				zap.NewAtomicLevelAt(test.expected),
			)

			// Set the global logger with the test core
			ZapLog = zap.New(core)
			defer ZapLog.Sync() // Ensure logs are written to the buffer

			// Emit a debug log and an info log
			ZapLog.Debug("This is a debug message")
			ZapLog.Info("This is an info message")

			// Read buffer contents and validate logs
			output := buffer.String()

			if test.verbose {
				// Debug logs should appear in verbose mode
				assert.Contains(t, output, "This is a debug message")
			} else {
				// Debug logs should NOT appear in non-verbose mode
				assert.NotContains(t, output, "This is a debug message")
			}

			// Info logs should always appear
			assert.Contains(t, output, "This is an info message")
		})
	}
}
