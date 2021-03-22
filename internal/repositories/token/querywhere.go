package token

import (
	"github.com/foolish15/shorten-url-service/internal/repositories"
	"gorm.io/gorm"
)

//WhereExpLessThan filter in array
type WhereExpLessThan struct {
	Exp uint
}

//DB implement interface
func (f WhereExpLessThan) DB(db repositories.DB) *gorm.DB {
	g := db.(*gorm.DB)
	return g.Where("`tokens`.`exp`<?", f.Exp)
}

//WhereSub filter in array
type WhereSub struct {
	Sub string
}

//DB implement interface
func (f WhereSub) DB(db repositories.DB) *gorm.DB {
	g := db.(*gorm.DB)
	return g.Where("`tokens`.`sub`=?", f.Sub)
}

//WhereCode filter in array
type WhereCode struct {
	Code string
}

//DB implement interface
func (f WhereCode) DB(db repositories.DB) *gorm.DB {
	g := db.(*gorm.DB)
	return g.Where("`tokens`.`code`=?", f.Code)
}
