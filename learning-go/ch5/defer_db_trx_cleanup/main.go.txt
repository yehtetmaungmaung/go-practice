package main

import (
	"context"
	"database/sql"
)

func DoSomeInserts(ctx context.Context, db *sql.DB, value1, value2 string) (err error) {
	trx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			err = trx.Commit()
		}
		if err != nil {
			trx.Rollback()
		}
	}()
	_, err = trx.ExecContext(ctx, "INSERT INTO FOO (val) values $1", value1)
	if err != nil {
		return err
	}
	return nil
}
