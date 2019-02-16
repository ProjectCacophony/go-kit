package migrations

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

// Logger provides a logger implementation to use with github.com/golang-migrate/migrate
type Logger struct {
	logger  *zap.Logger
	verbose bool
}

// NewLogger creates a new Logger
func NewLogger(logger *zap.Logger, verbose bool) *Logger {
	return &Logger{
		logger:  logger,
		verbose: verbose,
	}
}

// Printf logs the given message
func (l *Logger) Printf(format string, v ...interface{}) {
	l.logger.Info(
		strings.TrimSpace(fmt.Sprintf(format, v...)),
	)
}

// Verbose returns the verbosity configuration
func (l *Logger) Verbose() bool {
	return l.verbose
}
