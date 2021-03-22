package schemas

//UserAuthChannel enum user auth channel
type UserAuthChannel string

//define posible user auth channel
const (
	UserAuthChannelPassword UserAuthChannel = "password"
	UserAuthChannelSSO      UserAuthChannel = "sso"
	UserAuthChannelLine     UserAuthChannel = "line"
	UserAuthChannelFacebook UserAuthChannel = "facebook"
	UserAuthChannelGoogle   UserAuthChannel = "google"
)

//UserAuth schema user_auths
type UserAuth struct {
	ID            uint            `gorm:"primaryKey" json:"id"`
	UserID        uint            `gorm:"uniqueIndex:uidx_user_id_channel" json:"userId"`
	Channel       UserAuthChannel `gorm:"type:ENUM('password','line','facebook','google');uniqueIndex:uidx_user_id_channel" json:"channel"`
	ChannelID     string          `gorm:"type:varchar(255)" json:"channelId"`
	ChannelSecret string          `gorm:"type:varchar(255)" json:"-"`

	User User `json:"-"`
}
