package token

import (
	"github.com/foolish15/shorten-url-service/internal/repositories"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//ClauseLockForUpdate lock for update
type ClauseLockForUpdate struct{}

//DB implement interface
func (ClauseLockForUpdate) DB(db repositories.DB) *gorm.DB {
	g := db.(*gorm.DB)
	return g.Clauses(clause.Locking{Strength: "UPDATE"})
}
