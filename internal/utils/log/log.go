package log

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger       // Global Zap logger instance
var sugar *zap.SugaredLogger // Convenient sugared logger for fmt.Printf-style logging

// InitLogger initializes the global Zap logger based on the environment.
// This function should be called once at application startup.
func InitLogger(env string) {
	var config zap.Config
	if env == "production" {
		config = zap.NewProductionConfig() // Optimized for production logging
	} else {
		config = zap.NewDevelopmentConfig() // More human-readable for development
		// Add color to log levels in development for better readability
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Customize encoder configuration for consistent log output
	config.EncoderConfig.TimeKey = "timestamp"                   // Key for the timestamp field
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Use ISO8601 format for timestamps
	config.EncoderConfig.LevelKey = "level"                      // Key for the log level field
	config.EncoderConfig.MessageKey = "message"                  // Key for the log message field
	config.DisableCaller = true                                  // Disable caller info by default (can be enabled if needed)

	var err error
	logger, err = config.Build() // Build the logger
	if err != nil {
		fmt.Printf("Failed to build logger: %v\n", err)
		os.Exit(1) // Exit if logger cannot be initialized
	}
	sugar = logger.Sugar() // Create the sugared logger

	zap.ReplaceGlobals(logger) // Set as global logger for zap.L() usage
}

// Sync flushes any buffered log entries.
// This should be called before the application exits to ensure all logs are written.
func Sync() {
	if logger != nil {
		_ = logger.Sync() // Ignore error on sync
	}
}

// --- Wrapper functions for common log levels, allowing context passing ---

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

// --- Sugared Logger wrappers for convenience (fmt.Printf-style) ---
func Debugf(ctx context.Context, format string, args ...interface{}) {
	sugar.Debugf(format, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	sugar.Infof(format, args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	sugar.Warnf(format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	sugar.Errorf(format, args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	sugar.Fatalf(format, args...)
}
