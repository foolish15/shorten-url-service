package userauth

import (
	"github.com/foolish15/shorten-url-service/internal/repositories"
	"gorm.io/gorm"
)

//PreloadUser preload auth of user
type PreloadUser struct{}

//DB implement interface
func (f PreloadUser) DB(db repositories.DB) *gorm.DB {
	g := db.(*gorm.DB)
	return g.Preload("User")
}
