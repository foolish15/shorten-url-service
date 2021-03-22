package schemas

import (
	"time"

	"gorm.io/gorm"
)

type TokenSub string

const (
	TokenSubUser TokenSub = "user credential"
)

//Token token
type Token struct {
	ID        uint     `gorm:"primaryKey" json:"id"`
	Sub       TokenSub `gorm:"type:varchar(255)"`
	Exp       int64
	U         JSON           `gorm:"type:json" json:"u"`
	Code      *string        `gorm:"type:varchar(255);uniqueIndex" json:"-"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
