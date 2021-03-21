package schemas

import "time"

//AccessTransaction schema log
type AccessTransaction struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	LinkID         uint      `gorm:"index" json:"linkId"`
	LinkURL        string    `gorm:"type:varchar(255)" json:"targetUrl"`
	Browser        string    `gorm:"type:varchar(255);index" json:"browser"`
	BrowserVersion string    `json:"browserVersion"`
	OS             string    `gorm:"type:varchar(255);index" json:"os"`
	OSVersion      string    `json:"osVersion"`
	DeviceType     string    `gorm:"type:varchar(255);index" json:"deviceType"`
	UserAgent      string    `json:"userAgent"`
	CreatedAt      time.Time `json:"createAt"`
	UpdatedAt      time.Time `json:"updateAt"`
}
