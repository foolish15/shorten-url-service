package userauth

import (
	"github.com/foolish15/shorten-url-service/internal/repositories"
	"github.com/foolish15/shorten-url-service/internal/schemas"
	"gorm.io/gorm"
)

//WhereID filter in array
type WhereID struct {
	ID uint
}

//DB implement interface
func (f WhereID) DB(db repositories.DB) *gorm.DB {
	g := db.(*gorm.DB)
	return g.Where("`user_auths`.`id`=?", f.ID)
}

//WhereChannelID filter channelID
type WhereChannelID struct {
	ChannelID string
}

//DB implement interface
func (f WhereChannelID) DB(db repositories.DB) *gorm.DB {
	g := db.(*gorm.DB)
	return g.Where("`user_auths`.`channel_id`=?", f.ChannelID)
}

//WhereChannel filter channel
type WhereChannel struct {
	Channel schemas.UserAuthChannel
}

//DB implement interface
func (f WhereChannel) DB(db repositories.DB) *gorm.DB {
	g := db.(*gorm.DB)
	return g.Where("`user_auths`.`channel`=?", f.Channel)
}
