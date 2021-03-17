package block

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

//Create create block
func (g *Gorm) Create(bl *schemas.Block) error {
	db := g.DB
	err := db.Create(bl).Error
	if err != nil {
		return err
	}
	return nil
}

//Update create reward
func (g *Gorm) Update(bl *schemas.Block) error {
	db := g.DB
	err := db.Save(bl).Error
	if err != nil {
		return err
	}
	return nil
}

//Delete delete block
func (g *Gorm) Delete(t schemas.BlockType, v string) error {
	bl := schemas.Block{}
	bl, err := g.First(WhereType{Type: t}, WhereValue{Value: v})
	if err != nil {
		return err
	}
	db := g.DB
	err = db.Delete(&bl).Error
	if err != nil {
		return err
	}
	return nil
}

//Find find reward
func (g *Gorm) Find(queries ...repositories.Query) (bls []schemas.Block, err error) {
	db := g.Concat(g.DB, queries...)
	err = db.Find(&bls).Error
	if err != nil {
		return bls, err

	}
	return bls, nil
}

//First find reward
func (g *Gorm) First(queries ...repositories.Query) (bl schemas.Block, err error) {
	db := g.Concat(g.DB, queries...)
	err = db.First(&bl).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return bl, ErrNotFound
		}
		return bl, err

	}
	return bl, nil
}

//Count count shop
func (g *Gorm) Count(queries ...repositories.Query) (cnt int64, err error) {
	db := g.Concat(g.DB, queries...)
	err = db.Model(schemas.Block{}).Count(&cnt).Error
	if err != nil {
		return 0, err
	}
	return cnt, err
}

//TranslateOrderField translate order
func (g *Gorm) TranslateOrderField(input string) (output string) {
	list := map[string]string{
		"type":  "`blocks`.`type`",
		"value": "`blocks`.`value`",
	}

	output, ok := list[input]
	if !ok {
		return ""
	}
	return output
}

//New with argument
func (g *Gorm) New(idb ...interface{}) (Repository, error) {
	l := len(idb)
	if l > 1 {
		return nil, fmt.Errorf("cannot handle more than 1 argument")
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
		return nil, nil, fmt.Errorf("cannot start transaction: db is nil")
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
