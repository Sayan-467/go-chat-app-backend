package models

import "time"

type User struct {
	Id        uint		`gorm:"primaryKey" json:"id"`
	Username  string	`gorm:"unique;not null" json:"username"`
	Email     string	`gorm:"unique;not null" json:"email"`
	Password  string	`json:"-"`
	CreatedAt time.Time	`json:"created_at"`
}