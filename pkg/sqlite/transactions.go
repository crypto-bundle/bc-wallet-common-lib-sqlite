package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

var (
	ErrUnableGetTransactionFromContext = errors.New("unable get transaction from context")
	ErrNotInContextualTxStatement      = errors.New("unable to commit transaction statement - not in tx statement")
)

type transactionCtxKey string

//nolint:gochecknoglobals // it's ok
var transactionKey = transactionCtxKey("transaction")

// BeginTx ....
func (c *Connection) BeginTx() (*sqlx.Tx, error) {
	tx, err := c.Dbx.Beginx()
	if err != nil {
		return nil, c.e.ErrorOnly(err)
	}

	return tx, nil
}

// BeginTxWithRollbackOnError ....
func (c *Connection) BeginTxWithRollbackOnError(ctx context.Context,
	callback func(txStmtCtx context.Context) error,
) error {
	err := c.BeginReadCommittedTxRollbackOnError(ctx, callback)
	if err != nil {
		return c.e.ErrorNoWrap(err)
	}

	return nil
}

func (c *Connection) BeginReadCommittedTxRollbackOnError(ctx context.Context,
	callback func(txStmtCtx context.Context) error,
) error {
	txStmt, err := c.Dbx.Beginx()
	if err != nil {
		return c.e.ErrorOnly(err)
	}

	newCtx := context.WithValue(ctx, transactionKey, txStmt)

	err = callback(newCtx)
	if err != nil {
		rollbackErr := txStmt.Rollback()
		if rollbackErr != nil {
			return c.e.ErrorOnly(rollbackErr)
		}

		return c.e.ErrorOnly(err)
	}

	err = txStmt.Commit()
	if err != nil {
		return c.e.ErrorOnly(err)
	}

	return nil
}

func (c *Connection) BeginReadUncommittedTxRollbackOnError(ctx context.Context,
	callback func(txStmtCtx context.Context) error,
) error {
	txStmt, err := c.Dbx.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadUncommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return c.e.ErrorOnly(err)
	}

	newCtx := context.WithValue(ctx, transactionKey, txStmt)

	err = callback(newCtx)
	if err != nil {
		rollbackErr := txStmt.Rollback()
		if rollbackErr != nil {
			c.l.Warn("unable to rollback transaction, probably tx in pending status",
				slog.Any("error", rollbackErr))

			return c.e.ErrorOnly(rollbackErr)
		}

		return c.e.ErrorOnly(err)
	}

	err = txStmt.Commit()
	if err != nil {
		return c.e.ErrorOnly(err)
	}

	return nil
}

// BeginContextualTxStatement ....
func (c *Connection) BeginContextualTxStatement(ctx context.Context) (context.Context, error) {
	txStmt, err := c.Dbx.Beginx()
	if err != nil {
		return nil, c.e.ErrorOnly(err)
	}

	return context.WithValue(ctx, transactionKey, txStmt), nil
}

// CommitContextualTxStatement ....
func (c *Connection) CommitContextualTxStatement(ctx context.Context) error {
	tx, inTransaction := ctx.Value(transactionKey).(*sqlx.Tx)
	if !inTransaction {
		return c.e.ErrorOnly(ErrNotInContextualTxStatement)
	}

	err := tx.Commit()
	if err != nil {
		return c.e.ErrorOnly(err)
	}

	return nil
}

// RollbackContextualTxStatement ....
func (c *Connection) RollbackContextualTxStatement(ctx context.Context) error {
	tx, inTransaction := ctx.Value(transactionKey).(*sqlx.Tx)
	if !inTransaction {
		return c.e.ErrorOnly(ErrNotInContextualTxStatement)
	}

	err := tx.Rollback()
	if err != nil {
		return c.e.ErrorOnly(err)
	}

	return nil
}

func (c *Connection) TryWithTransaction(ctx context.Context, sqlExecutionFunc func(stmt sqlx.Ext) error) error {
	stmt := sqlx.Ext(c.Dbx)

	tx, inTransaction := ctx.Value(transactionKey).(*sqlx.Tx)
	if inTransaction {
		stmt = tx
	}

	return sqlExecutionFunc(stmt)
}

func (c *Connection) MustWithTransaction(ctx context.Context, sqlInTxExecutionFunc func(stmt *sqlx.Tx) error) error {
	tx, inTransaction := ctx.Value(transactionKey).(*sqlx.Tx)
	if inTransaction {
		return sqlInTxExecutionFunc(tx)
	}

	return c.e.ErrorOnly(ErrUnableGetTransactionFromContext)
}
