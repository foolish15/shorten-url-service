package user

import (
	"github.com/foolish15/shorten-url-service/internal/repositories"
	"gorm.io/gorm"
)

//PreloadUserAuth preload auth of user
type PreloadUserAuth struct{}

//DB implement interface
func (f PreloadUserAuth) DB(db repositories.DB) *gorm.DB {
	g := db.(*gorm.DB)
	return g.Preload("UserAuths")
}
