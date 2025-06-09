package sqlite

import (
	"log/slog"
	"time"
)

//nolint:interfacebloat //it's ok here, we need it we must use it as one big interface
type errorFormatterService interface {
	// ErrorNoWrap function for pseudo-wrap error, must be used in case of linter warnings...
	ErrorNoWrap(err error) error
	// ErrNoWrap same with ErrorNoWrap function, just alias for ErrorNoWrap, just short function name...
	ErrNoWrap(err error) error
	ErrorOnly(err error, details ...string) error
	Error(err error, details ...string) error
	Errorf(err error, format string, args ...interface{}) error
	NewError(details ...string) error
	NewErrorf(format string, args ...interface{}) error
}

type loggerBuilderService interface {
	NewSlogLoggerEntry(fields ...any) *slog.Logger
	NewSlogNamedLoggerEntry(named string, fields ...any) *slog.Logger
	NewSlogLoggerEntryWithFields(fields ...slog.Attr) *slog.Logger
}

type configSQLiteParametersService interface {
	GetDatabaseDSN() string
	GetDBName() string
	GetDBUser() string
	GetDBPassword() string
	GetDBConnectTimeOut() time.Duration
	GetSQLiteDBFilePath() string
	GetConnectionRetryCount() uint8
	GetConnectionRetryTimeout() time.Duration
}
