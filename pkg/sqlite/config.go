package sqlite

import (
	"fmt"
	"time"
)

var _ configSQLiteParametersService = (*SQLiteConfig)(nil)

type SQLiteConfig struct {
	DBFilePath string `envconfig:"SQLITE_DATABASE_FILE_PATH" default:"/var/lib/application/db.sqlite"`
	DBName     string `envconfig:"SQLITE_DATABASE_DATABASE_NAME" default:"ca-api-gateway"`
	DBUsername string `envconfig:"SQLITE_DATABASE_USERNAME"`
	DBPassword string `envconfig:"SQLITE_DATABASE_PASSWORD"`
	// DBConnectTimeOut is the timeout in millisecond to connect between connection tries
	DBConnectTimeOut uint16 `envconfig:"SQLITE_CONNECTION_RETRY_TIMEOUT" default:"5000"`
	// DBConnectRetryCount is the maximum number of reconnection tries. If 0 - infinite loop
	DBConnectRetryCount uint8 `envconfig:"SQLITE_CONNECTION_RETRY_COUNT" default:"0"`

	// calculated parameters
	retryTimeOut time.Duration
}

func (c *SQLiteConfig) Prepare() error {
	c.retryTimeOut = time.Duration(c.GetDBConnectTimeOut()) * time.Millisecond

	return nil
}

func (c *SQLiteConfig) GetDatabaseDSN() string {
	return fmt.Sprintf("file:%s?_mutex=no&mode=rwc&_txlock=immediate", c.DBFilePath)
}

func (c *SQLiteConfig) GetSQLiteDBFilePath() string {
	return c.DBFilePath
}

func (c *SQLiteConfig) GetDBName() string {
	return c.DBName
}

func (c *SQLiteConfig) GetDBUser() string {
	return c.DBUsername
}

func (c *SQLiteConfig) GetDBPassword() string {
	return c.DBPassword
}

func (c *SQLiteConfig) GetDBRetryCount() uint8 {
	return c.DBConnectRetryCount
}

func (c *SQLiteConfig) GetDBConnectTimeOut() uint16 {
	return c.DBConnectTimeOut
}

func (c *SQLiteConfig) GetConnectionRetryCount() uint8 {
	return c.DBConnectRetryCount
}

func (c *SQLiteConfig) GetConnectionRetryTimeout() time.Duration {
	return c.retryTimeOut
}
