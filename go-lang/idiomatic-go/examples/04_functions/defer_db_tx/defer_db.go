// defer + именованная возвращаемая ошибка = commit/rollback транзакции в одном месте.
// Отложенная функция читает err ПОСЛЕ return и решает, фиксировать или откатывать,
// а также может подменить возвращаемую ошибку (ошибкой самого Commit).
package defer_db

import (
	"context"
	"database/sql"
)

func DoSomeInserts(ctx context.Context, db *sql.DB, value1, value2 string) (err error) {
	// ИМЕНОВАННАЯ ошибка — без нее отложенная функция не увидела бы результат
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			// defer читает err уже ПОСЛЕ return и решает: commit или rollback
			err = tx.Commit()
			// и может подменить возвращаемую ошибку — например ошибкой самого Commit
		}
		if err != nil {
			tx.Rollback()
		}
	}()
	_, err = tx.ExecContext(ctx, "INSERT INTO FOO (val) values $1", value1)
	if err != nil {
		return err
	}
	// use tx to do more database inserts here
	return nil
}
