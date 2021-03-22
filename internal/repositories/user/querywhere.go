package user

import (
	"github.com/foolish15/shorten-url-service/internal/repositories"
	"gorm.io/gorm"
)

//WhereID filter in array
type WhereID struct {
	ID uint
}

//DB implement interface
func (f WhereID) DB(db repositories.DB) *gorm.DB {
	g := db.(*gorm.DB)
	return g.Where("`users`.`id`=?", f.ID)
}

//WhereEmail filter in array
type WhereEmail struct {
	Email string
}

//DB implement interface
func (f WhereEmail) DB(db repositories.DB) *gorm.DB {
	g := db.(*gorm.DB)
	return g.Where("`users`.`email`=?", f.Email)
}

//WherePhone filter in array
type WherePhone struct {
	Phone string
}

//DB implement interface
func (f WherePhone) DB(db repositories.DB) *gorm.DB {
	g := db.(*gorm.DB)
	return g.Where("`users`.`phone`=?", f.Phone)
}
