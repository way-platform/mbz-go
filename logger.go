package mbz

import "log/slog"

// Logger is a leveled logger interface.
type Logger interface {
	Debug(msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
}

type noopLogger struct{}

var _ Logger = noopLogger{}

func (noopLogger) Debug(msg string, keysAndValues ...any) {}
func (noopLogger) Info(msg string, keysAndValues ...any)  {}
func (noopLogger) Warn(msg string, keysAndValues ...any)  {}
func (noopLogger) Error(msg string, keysAndValues ...any) {}

type slogLogger struct {
	logger *slog.Logger
}

func (l slogLogger) Debug(msg string, keysAndValues ...any) {
	l.logger.Debug(msg, keysAndValues...)
}

func (l slogLogger) Info(msg string, keysAndValues ...any) {
	l.logger.Info(msg, keysAndValues...)
}

func (l slogLogger) Warn(msg string, keysAndValues ...any) {
	l.logger.Warn(msg, keysAndValues...)
}

func (l slogLogger) Error(msg string, keysAndValues ...any) {
	l.logger.Error(msg, keysAndValues...)
}
