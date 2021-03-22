package userauth

import (
	"database/sql"
	"fmt"

	"github.com/foolish15/shorten-url-service/internal/repositories"
	"github.com/foolish15/shorten-url-service/internal/schemas"
	"gorm.io/gorm"
)

//Gorm type repo for gorm
type Gorm struct {
	*repositories.BaseGormRepository
}

//New new object
func New(db *gorm.DB) *Gorm {
	return &Gorm{
		BaseGormRepository: &repositories.BaseGormRepository{
			DB: db,
		},
	}
}

//Create create token
func (g *Gorm) Create(usrA *schemas.UserAuth) error {
	db := g.DB
	err := db.Create(usrA).Error
	if err != nil {
		return err
	}
	return nil
}

//Update create reward
func (g *Gorm) Update(usrA *schemas.UserAuth) error {
	db := g.DB
	err := db.Save(usrA).Error
	if err != nil {
		return err
	}
	return nil
}

//Delete delete token
func (g *Gorm) Delete(ID uint) error {
	usrA, err := g.First(WhereID{ID: ID})
	if err != nil {
		return err
	}
	db := g.DB
	err = db.Delete(&usrA).Error
	if err != nil {
		return err
	}
	return nil
}

//Deletes delete token
func (g *Gorm) Deletes(queries ...repositories.Query) error {
	db := g.Concat(g.DB, queries...)
	err := db.Delete(schemas.UserAuth{}).Error
	if err != nil {
		return err
	}
	return nil
}

//Find find reward
func (g Gorm) Find(queries ...repositories.Query) (usrAs []schemas.UserAuth, err error) {
	db := g.Concat(g.DB, queries...)
	err = db.Find(&usrAs).Error
	if err != nil {
		return usrAs, err

	}
	return usrAs, nil
}

//First find reward
func (g Gorm) First(queries ...repositories.Query) (usrA schemas.UserAuth, err error) {
	db := g.Concat(g.DB, queries...)
	err = db.First(&usrA).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return usrA, ErrNotFound
		}
		return usrA, err

	}
	return usrA, nil
}

//Count count shop
func (g Gorm) Count(queries ...repositories.Query) (cnt int64, err error) {
	db := g.Concat(g.DB, queries...)
	err = db.Model(schemas.UserAuth{}).Count(&cnt).Error
	if err != nil {
		return 0, err
	}
	return cnt, err
}

//New with argument
func (g *Gorm) New(idb ...interface{}) (Repository, error) {
	l := len(idb)
	if l > 1 {
		return nil, fmt.Errorf("Cannot handle more than 1 argument")
	}
	newRepo := New(g.DB)
	if l == 1 {
		err := newRepo.SetDB(idb[0])
		return newRepo, err
	}

	return newRepo, nil
}

//StartTransaction start transaction
func (g *Gorm) StartTransaction() (tx *sql.Tx, rp Repository, err error) {
	db, err := g.DB.DB()
	if err != nil {
		return nil, nil, err
	}
	if db == nil {
		return nil, nil, fmt.Errorf("Cannot start transaction: db is nil")
	}
	gsess := g.DB.Session(&gorm.Session{Context: g.DB.Statement.Context})
	var opt *sql.TxOptions

	beginner, ok := gsess.Statement.ConnPool.(*sql.DB)
	if !ok {
		return nil, nil, gorm.ErrInvalidTransaction
	}

	tx, err = beginner.BeginTx(gsess.Statement.Context, opt)
	if err != nil {
		return nil, nil, err
	}

	txRepo, err := g.New(tx)
	if err != nil {
		return nil, nil, err
	}

	return tx, txRepo, nil
}
