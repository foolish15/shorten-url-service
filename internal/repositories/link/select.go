package link

import (
	"github.com/foolish15/shorten-url-service/internal/repositories"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//SelectUnscope lock select for update
type SelectUnscope struct{}

//DB implement interface
func (w SelectUnscope) DB(db repositories.DB) *gorm.DB {
	g := db.(*gorm.DB)
	return g.Unscoped()
}

//SelectForUpdate lock select for update
type SelectForUpdate struct{}

//DB implement interface
func (w SelectForUpdate) DB(db repositories.DB) *gorm.DB {
	g := db.(*gorm.DB)
	return g.Clauses(clause.Locking{Strength: "UPDATE"})
}
