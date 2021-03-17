package block

import (
	"github.com/foolish15/shorten-url-service/internal/repositories"
	"github.com/foolish15/shorten-url-service/internal/schemas"
	"gorm.io/gorm"
)

//WhereType filter brandID
type WhereType struct {
	Type schemas.BlockType
}

//DB implement interface
func (w WhereType) DB(db repositories.DB) *gorm.DB {
	g := db.(*gorm.DB)
	return g.Where("`blocks`.`type` = ?", w.Type)
}

//WhereValue filter brandID
type WhereValue struct {
	Value string
}

//DB implement interface
func (w WhereValue) DB(db repositories.DB) *gorm.DB {
	g := db.(*gorm.DB)
	return g.Where("`blocks`.`value` = ?", w.Value)
}
