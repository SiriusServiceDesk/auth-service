package models

type User struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	Surname    string  `json:"surname"`
	SecondName *string `json:"SecondName"`
	Password   string  `json:"-"`
	Email      string  `json:"email"`
}
