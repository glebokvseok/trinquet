package psql

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	errorMessage = "unhandled panic occurred while executing postgresql request: %v"
)

func PanicHandler(returnErr *error) {
	if r := recover(); r != nil {
		if returnErr != nil {
			*returnErr = errors.Errorf(errorMessage, r)
		}
	}
}

func TransactionFinalizer(tx *gorm.DB, returnErr *error) {
	if r := recover(); r != nil {
		if tx != nil {
			tx.Rollback()
		}

		if returnErr != nil {
			*returnErr = errors.Errorf(errorMessage, r)
		}

		return
	}

	if returnErr != nil && *returnErr != nil && tx != nil {
		tx.Rollback()
	}
}
