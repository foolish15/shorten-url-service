package token

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
func (g *Gorm) Create(tkn *schemas.Token) error {
	db := g.DB
	err := db.Create(tkn).Error
	if err != nil {
		return err
	}
	return nil
}

//Update create reward
func (g *Gorm) Update(tkn *schemas.Token) error {
	db := g.DB
	err := db.Save(tkn).Error
	if err != nil {
		return err
	}
	return nil
}

//Delete delete token
func (g *Gorm) Delete(ID uint) error {
	tkn := schemas.Token{}
	db := g.DB
	err := db.First(&tkn, "id = ?", ID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}
	err = db.Delete(&tkn).Error
	if err != nil {
		return err
	}
	return nil
}

//Deletes delete token
func (g *Gorm) Deletes(queries ...repositories.Query) error {
	db := g.Concat(g.DB, queries...)
	err := db.Delete(schemas.Token{}).Error
	if err != nil {
		return err
	}
	return nil
}

//Find find reward
func (g Gorm) Find(queries ...repositories.Query) (tkns []schemas.Token, err error) {
	db := g.Concat(g.DB, queries...)
	err = db.Find(&tkns).Error
	if err != nil {
		return tkns, err

	}
	return tkns, nil
}

//First find reward
func (g Gorm) First(queries ...repositories.Query) (tkn schemas.Token, err error) {
	db := g.Concat(g.DB, queries...)
	err = db.First(&tkn).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return tkn, ErrNotFound
		}
		return tkn, err

	}
	return tkn, nil
}

//Count count shop
func (g Gorm) Count(queries ...repositories.Query) (cnt int64, err error) {
	db := g.Concat(g.DB, queries...)
	err = db.Model(schemas.Token{}).Count(&cnt).Error
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
	newRepo := &Gorm{
		BaseGormRepository: &repositories.BaseGormRepository{
			DB: g.DB,
		},
	}
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
