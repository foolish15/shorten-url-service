package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"gorm.io/gorm"
)

//SQLDb interface validate DB
type SQLDb interface {
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

//SQLTx interface validate tx
type SQLTx interface {
	Commit() error
	Rollback() error
}

//SetDBInterface for db repository
type SetDBInterface interface {
	SetDB(idb interface{}) error
}

// DB interface
type DB interface{}

// Query type interface
type Query interface {
	DB(DB) *gorm.DB
}

//BaseGormRepository base gorm repo
type BaseGormRepository struct {
	DB *gorm.DB
}

//Concat concat where
func (r *BaseGormRepository) Concat(db *gorm.DB, queries ...Query) *gorm.DB {
	for _, query := range queries {
		db = query.DB(db)
	}
	return db
}

//SetDB set DB
func (r *BaseGormRepository) SetDB(idb interface{}) error {
	switch idb.(type) {
	case *gorm.DB:
		r.DB = idb.(*gorm.DB)
		return nil
	case *sql.DB:
		gsess := r.DB.Session(&gorm.Session{Context: r.DB.Statement.Context})
		sqlDB := idb.(*sql.DB)
		gsess.Statement.ConnPool = sqlDB
		return r.SetDB(gsess)
	case *sql.Tx:
		gsess := r.DB.Session(&gorm.Session{Context: r.DB.Statement.Context})
		sqlTx := idb.(*sql.Tx)
		gsess.Statement.ConnPool = sqlTx
		return r.SetDB(gsess)
	}
	return fmt.Errorf("Cannot set db with idb argument type[%T]", idb)
}
