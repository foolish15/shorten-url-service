package link

import (
	"github.com/foolish15/shorten-url-service/internal/repositories"
	"gorm.io/gorm"
)

//WhereID filter brandID
type WhereID struct {
	ID uint
}

//DB implement interface
func (w WhereID) DB(db repositories.DB) *gorm.DB {
	g := db.(*gorm.DB)
	return g.Where("`links`.`id` = ?", w.ID)
}
