package helpers

import (
	"gorm.io/gorm"
)

func runInTx(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// execute user function inside transaction
		return fn(tx)
	})
}

// func runInTx(db *sql.DB, fn func(tx *sql.Tx) error) error {
// 	tx, err := db.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	err = fn(tx)
// 	if err == nil {
// 		return tx.Commit()
// 	}

// 	rollbackErr := tx.Rollback()
// 	if rollbackErr != nil {
// 		return errors.Join(err, rollbackErr)
// 	}

// 	return err
// }
