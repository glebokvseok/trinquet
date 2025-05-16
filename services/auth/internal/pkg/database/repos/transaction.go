package repos

import (
	"context"
	txmgr "github.com/move-mates/trinquet/library/database/pkg/psql/managers"
	"gorm.io/gorm"
)

func ProvideTransactionManager(
	db *gorm.DB,
) txmgr.TransactionManager {
	return txmgr.NewTransactionManager(db, txKey{})
}

func getTransaction(ctx context.Context) (*gorm.DB, error) {
	return txmgr.GetTransaction(ctx, txKey{})
}

type txKey struct{}
