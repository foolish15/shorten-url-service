package accesstransaction

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

//Create create access transaction
func (g *Gorm) Create(acc *schemas.AccessTransaction) error {
	db := g.DB
	err := db.Create(acc).Error
	if err != nil {
		return err
	}
	return nil
}

//Update update access transaction
func (g *Gorm) Update(acc *schemas.AccessTransaction) error {
	db := g.DB
	err := db.Save(acc).Error
	if err != nil {
		return err
	}
	return nil
}

//Delete delete access transaction
func (g *Gorm) Delete(id uint) error {
	acc := schemas.AccessTransaction{}
	acc, err := g.First(WhereID{ID: id})
	if err != nil {
		return err
	}
	db := g.DB
	err = db.Delete(&acc).Error
	if err != nil {
		return err
	}
	return nil
}

//Find find access transaction
func (g *Gorm) Find(queries ...repositories.Query) (accs []schemas.AccessTransaction, err error) {
	db := g.Concat(g.DB, queries...)
	err = db.Find(&accs).Error
	if err != nil {
		return accs, err

	}
	return accs, nil
}

//First find access transaction
func (g *Gorm) First(queries ...repositories.Query) (acc schemas.AccessTransaction, err error) {
	db := g.Concat(g.DB, queries...)
	err = db.First(&acc).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return acc, ErrNotFound
		}
		return acc, err

	}
	return acc, nil
}

//Count count access transaction
func (g *Gorm) Count(queries ...repositories.Query) (cnt int64, err error) {
	db := g.Concat(g.DB, queries...)
	err = db.Model(schemas.AccessTransaction{}).Count(&cnt).Error
	if err != nil {
		return 0, err
	}
	return cnt, err
}

//TranslateOrderField translate order
func (g *Gorm) TranslateOrderField(input string) (output string) {
	list := map[string]string{
		"link": "`links`.`link`",
		"code": "`links`.`code`",
		"exp":  "`links`.`expire`",
		"hit":  "`links`.`hit`",
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
