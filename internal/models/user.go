package models

import "time"

type userRole string

const (
	Admin         userRole = "Админ"
	TechDep       userRole = "Технический отдел"
	MethodicalDep userRole = "Методический отдел"
	StudyOffice   userRole = "Учебный офис"
	Hotel         userRole = "Гостиница"
	Educators     userRole = "Воспитательный отдел"
	UserR         userRole = "user"
)

type User struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Surname    string   `json:"surname"`
	SecondName *string  `json:"second_name"`
	Password   string   `json:"-"`
	Email      string   `json:"email" gorm:"unique:true"`
	IsVerified bool     `json:"is_verified"`
	TelegramId string   `json:"telegram_id"`
	Role       userRole `json:"role"`

	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (u *User) IsAdmin() bool {
	return u.Role == Admin
}
