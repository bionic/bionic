package provider

import (
	"errors"
	"gorm.io/gorm"
)

var (
	ErrTxAlreadyInProgress = errors.New("tx already in progress")
	ErrNoTxInProgress      = errors.New("no tx in progress")
)

type Database interface {
	DB() *gorm.DB
	BeginTx() error
	CommitTx() error
	RollbackTx() error
}

type database struct {
	db *gorm.DB
	tx *gorm.DB
}

func NewDatabase(db *gorm.DB) Database {
	return &database{
		db: db,
	}
}

func (db *database) DB() *gorm.DB {
	if db.tx != nil {
		return db.tx
	}

	return db.db
}

func (db *database) BeginTx() error {
	if db.tx != nil {
		return ErrTxAlreadyInProgress
	}

	tx := db.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	db.tx = tx

	return nil
}

func (db *database) CommitTx() error {
	if db.tx == nil {
		return ErrNoTxInProgress
	}

	commit := db.tx.Commit()
	db.tx = nil

	if err := commit.Error; err != nil {
		return err
	}

	return nil
}

func (db *database) RollbackTx() error {
	if db.tx == nil {
		return nil
	}

	rollback := db.tx.Rollback()
	db.tx = nil

	if err := rollback.Error; err != nil {
		return err
	}

	return nil
}
