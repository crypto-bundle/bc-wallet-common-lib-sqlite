package sqlite

import "log/slog"

var _ loggerBuilderService = (*defaultLogBuilder)(nil)

type defaultLogBuilder struct {
}

func (s *defaultLogBuilder) NewSlogLoggerEntry(fields ...any) *slog.Logger {
	return slog.Default().With(fields...)
}
func (s *defaultLogBuilder) NewSlogNamedLoggerEntry(named string, fields ...any) *slog.Logger {
	return slog.Default().WithGroup(named).With(fields...)
}
func (s *defaultLogBuilder) NewSlogLoggerEntryWithFields(fields ...slog.Attr) *slog.Logger {
	return slog.Default().With(fields)
}

func NewDefaultSQLiteLoggerBuilder() *defaultLogBuilder {
	return &defaultLogBuilder{}
}
