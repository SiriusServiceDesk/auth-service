package repository

import (
	"errors"
	"github.com/SiriusServiceDesk/auth-service/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

func (u *UserRepositoryImpl) seeds() error {
	nowTime := time.Now()
	adminPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	users := []models.User{
		{
			Id:         uuid.NewString(),
			Name:       "админ",
			Surname:    "админов",
			SecondName: nil,
			Password:   string(adminPassword),
			Email:      "admin@siriusdesk.ru",
			IsVerified: true,
			Role:       models.Admin,
			CreatedAt:  &nowTime,
		},
		{
			Id:         uuid.NewString(),
			Name:       "второй",
			Surname:    "админ",
			SecondName: nil,
			Password:   string(adminPassword),
			Email:      "admin2@siriusdesk.ru",
			IsVerified: true,
			Role:       models.Admin,
			CreatedAt:  &nowTime,
		},
	}
	for _, user := range users {
		_, err := u.GetUserByEmail(user.Email)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u.CreateUser(&user)
		}
	}
	return nil
}
