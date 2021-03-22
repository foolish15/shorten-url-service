package schemas

import (
	"time"

	"gorm.io/gorm"
)

//User table user
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      *string        `gorm:"type:varchar(255)" json:"name"`
	Email     *string        `gorm:"type:varchar(255)" json:"email"`
	Phone     *string        `gorm:"type:varchar(255)" json:"phone"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserAuths []UserAuth `json:"-"`
}
