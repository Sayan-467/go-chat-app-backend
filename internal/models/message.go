package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	Id        uint           `gorm:"primaryKey" json:"id"`
	Sender    string         `gorm:"sender"`
	Receiver  string         `gorm:"receiver"`
	Room      string         `gorm:"room"`
	Content   string         `gorm:"content"`
	CreatedAt time.Time      `gorm:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}