package sqlite

import (
	"database/sql"
	"errors"
)

func SetDefaultErrorFormatter(errFmtSvc errorFormatterService) {
	defaultErrorsFmtSvc = errFmtSvc
}

func EmptyOrError(err error, errorMessage string) error {
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	return ErrorOnly(err, errorMessage)
}
