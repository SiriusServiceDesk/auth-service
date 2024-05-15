package auth

import "github.com/SiriusServiceDesk/auth-service/internal/models"

func MappingUserModelToResponse(user *models.User) UserResponse {
	return UserResponse{
		Id:         user.Id,
		Name:       user.Name,
		Surname:    user.Surname,
		SecondName: user.SecondName,
		Email:      user.Email,
		TelegramId: user.TelegramId,
	}
}
