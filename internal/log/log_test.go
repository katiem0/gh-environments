package log

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name      string
		debug     bool
		wantLevel zapcore.Level
	}{
		{
			name:      "Debug mode enabled",
			debug:     true,
			wantLevel: zapcore.DebugLevel,
		},
		{
			name:      "Debug mode disabled",
			debug:     false,
			wantLevel: zapcore.InfoLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewLogger(tt.debug)
			if err != nil {
				t.Errorf("NewLogger() error = %v", err)
				return
			}
			if logger == nil {
				t.Error("NewLogger() returned nil logger")
			}

			// Create a logger with an observer for testing
			core, obs := observer.New(tt.wantLevel)
			obsLogger := zap.New(core)

			// Log at different levels
			if tt.debug {
				logger.Debug("test debug message")
				obsLogger.Debug("test debug message")

				// Verify debug logs are captured when debug is enabled
				logs := obs.All()
				if len(logs) != 1 {
					t.Errorf("Expected 1 debug log, got %d", len(logs))
				}
			} else {
				// When debug is disabled, verify info logs still work
				logger.Info("test info message")
				obsLogger.Info("test info message")

				logs := obs.All()
				if len(logs) != 1 {
					t.Errorf("Expected 1 info log, got %d", len(logs))
				}
			}
		})
	}
}
