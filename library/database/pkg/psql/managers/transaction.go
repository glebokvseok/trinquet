package managers

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TransactionManager interface {
	BeginTransaction(ctx context.Context) (context.Context, error)
	CommitTransaction(ctx context.Context) error
	RollbackTransaction(ctx context.Context)
	FinalizeTransaction(ctx context.Context, err *error)
}

type transactionManager struct {
	db    *gorm.DB
	txKey any
}

func NewTransactionManager(
	db *gorm.DB,
	txKey any,
) TransactionManager {
	return &transactionManager{
		db:    db,
		txKey: txKey,
	}
}

func (mgr *transactionManager) BeginTransaction(ctx context.Context) (context.Context, error) {
	if _, err := GetTransaction(ctx, mgr.txKey); err == nil {
		return ctx, errors.New("transaction already exists in context")
	}

	tx := mgr.db.Begin()
	if tx.Error != nil {
		return ctx, errors.WithStack(tx.Error)
	}

	return context.WithValue(ctx, mgr.txKey, tx), nil
}

func (mgr *transactionManager) CommitTransaction(ctx context.Context) error {
	tx, err := GetTransaction(ctx, mgr.txKey)
	if err != nil {
		return err
	}

	if err = tx.Commit().Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (mgr *transactionManager) RollbackTransaction(ctx context.Context) {
	tx, err := GetTransaction(ctx, mgr.txKey)
	if err != nil {
		return
	}

	tx.Rollback()
}

func (mgr *transactionManager) FinalizeTransaction(ctx context.Context, err *error) {
	if r := recover(); r != nil {
		mgr.RollbackTransaction(ctx)

		if err != nil {
			*err = errors.Errorf("unhandeled panic occurred while executing postgresql repository transaction: %v", r)
		}

		return
	}

	if err != nil && *err != nil {
		mgr.RollbackTransaction(ctx)
	}
}

func GetTransaction(ctx context.Context, txKey any) (*gorm.DB, error) {
	tx, ok := ctx.Value(txKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("no transaction found in context")
	}

	return tx, nil
}
