package schemas

import (
	"fmt"
	"os"
	"time"

	"gorm.io/gorm"
)

type Link struct {
	ID        uint           `gorm:"primaryKey" json:"-"`
	Link      string         `gorm:"type:text" json:"link"`
	Expire    int64          `json:"expire"`
	Code      string         `gorm:"type:varchar(255);uniqueIndex" json:"-"`
	Hit       uint           `gorm:"default:0" json:"hit"`
	AccessURL string         `gorm:"-" json:"accessUrl"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (l *Link) AfterFind(tx *gorm.DB) (err error) {
	l.AccessURL = fmt.Sprintf("%s/%s", os.Getenv("APP_URL"), l.Code)
	return
}

func (l *Link) AfterUpdate(tx *gorm.DB) (err error) {
	l.AccessURL = fmt.Sprintf("%s/link/%s", os.Getenv("APP_URL"), l.Code)
	return
}
