package models

import "time"

type Reason struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Content   string `gorm:"not null"`
	UserID    uint   `gorm:"not null"`
	User      User   `gorm:"foreignKey:UserID"`
	CompanyName string  `gorm:"not null"`
	CompanyUrl  string  `gorm:"not null"`
	CreatedAt time.Time
}
