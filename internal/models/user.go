package models

import "time"

type User struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	Surname    string  `json:"surname"`
	SecondName *string `json:"second_name"`
	Password   string  `json:"-"`
	Email      string  `json:"email" gorm:"unique:true"`
	IsVerified bool    `json:"is_verified"`
	TelegramId string  `json:"telegram_id"`

	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
