package sqlite

import (
	"context"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

// Connection struct to store and manipulate postgres database connection...
type Connection struct {
	l *slog.Logger
	e errorFormatterService
	c configSQLiteParametersService

	Dbx *sqlx.DB
}

func (c *Connection) IsHealed(ctx context.Context) bool {
	err := c.Dbx.PingContext(ctx)
	if err != nil {
		return false
	}

	err = c.checkConnectionByQuery(c.Dbx)

	return err == nil
}

func (c *Connection) checkConnectionByQuery(dbx *sqlx.DB) error {
	rows, err := dbx.Query("SELECT 1")
	if err != nil {
		return c.e.ErrorOnly(err)
	}

	defer func() {
		_ = rows.Close()
	}()

	err = rows.Err()
	if err != nil {
		return c.e.ErrorOnly(err)
	}

	return nil
}

// Connect to sqlite database...
func (c *Connection) Connect() (*Connection, error) {
	retryDecValue := uint8(1)
	retryCount := c.c.GetConnectionRetryCount()

	if retryCount == 0 {
		retryDecValue = 0
		retryCount = 1
	}

	try := 0

	var err error

	for i := retryCount; i != 0; i -= retryDecValue {
		dbx, loopErr := c.tryConnect()
		if loopErr != nil {
			c.l.Error("unable to connect to database", slog.Any("error", loopErr),
				slog.Int(ConnectionRetryCountTag, try))

			err = loopErr

			time.Sleep(c.c.GetConnectionRetryTimeout())

			continue
		}

		c.Dbx = dbx

		return c, nil
	}

	if err != nil {
		return nil, c.e.ErrorOnly(err)
	}

	return c, nil
}

func (c *Connection) tryConnect() (*sqlx.DB, error) {
	dbx, err := sqlx.Connect("sqlite3", c.c.GetSQLiteDBFilePath())
	if err != nil {
		c.l.Error("unable to connect sqlite database", slog.Any("error", err))
	}

	err = dbx.Ping()
	if err != nil {
		return nil, c.e.Error(err)
	}

	err = c.checkConnectionByQuery(dbx)
	if err != nil {
		return nil, c.e.Error(err)
	}

	return dbx, nil
}

func (c *Connection) Close() error {
	err := c.Dbx.Close()
	if err != nil {
		return c.e.ErrorOnly(err)
	}

	return nil
}

// NewConnection to postgres db...
func NewConnection(logFactorySvc loggerBuilderService,
	errFormatterSvc errorFormatterService,
	cfgSvc configSQLiteParametersService,
) *Connection {
	conn := &Connection{
		e:   errFormatterSvc,
		l:   logFactorySvc.NewSlogNamedLoggerEntry("lib-postgres"),
		c:   cfgSvc,
		Dbx: nil,
	}

	return conn
}

func NewConnectionWithDefaults(dbConn *sqlx.DB,
	logFactorySvc loggerBuilderService,
	errFmtSvc errorFormatterService,
) *Connection {
	return &Connection{
		l:   logFactorySvc.NewSlogNamedLoggerEntry("lib-postgres"),
		e:   errFmtSvc,
		c:   nil,
		Dbx: dbConn,
	}
}
