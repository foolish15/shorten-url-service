package link

import (
	"github.com/foolish15/shorten-url-service/internal/repositories"
	"gorm.io/gorm"
)

//LimitOffset type for limit query
type LimitOffset struct {
	Limit, Offset int
}

//DB implement interface
func (l LimitOffset) DB(db repositories.DB) *gorm.DB {
	g := db.(*gorm.DB)
	return g.Limit(l.Limit).Offset(l.Offset)
}

//LimitPage type for limit query
type LimitPage struct {
	Limit, Page int
}

//DB implement interface
func (l LimitPage) DB(db repositories.DB) *gorm.DB {
	g := db.(*gorm.DB)
	offset := (l.Page - 1) * l.Limit
	return g.Limit(l.Limit).Offset(offset)
}
