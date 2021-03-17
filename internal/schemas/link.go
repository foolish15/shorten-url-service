package schemas

import (
	"time"

	"gorm.io/gorm"
)

type Link struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Link      string         `gorm:"type:text" json:"link"`
	Expire    int64          `json:"expire"`
	Code      string         `gorm:"type:varchar(255);uniqueIndex" json:"code"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
