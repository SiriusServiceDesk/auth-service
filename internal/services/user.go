package services

import (
	"errors"
	"github.com/SiriusServiceDesk/auth-service/internal/config"
	"github.com/SiriusServiceDesk/auth-service/internal/models"
	"github.com/SiriusServiceDesk/auth-service/internal/repository"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserById(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUsers() ([]*models.User, error)
	UpdateUser(id string, user *models.User) (*models.User, error)
	DeleteUser(id string) error
	ComparePassword(userPassword, inputPassword string) error
	GenerateToken(user *models.User) (string, error)
	HashingPassword(password string) ([]byte, error)
}

func (u UserServiceImpl) GetUserByEmail(email string) (*models.User, error) {
	return u.repos.GetUserByEmail(email)
}

func (u UserServiceImpl) CreateUser(user *models.User) (*models.User, error) {
	user.Id = uuid.NewString()
	return u.repos.CreateUser(user)
}

func (u UserServiceImpl) GetUserById(id string) (*models.User, error) {
	return u.repos.GetUserById(id)
}

func (u UserServiceImpl) GetUsers() ([]*models.User, error) {
	return u.repos.GetUsers()
}

func (u UserServiceImpl) UpdateUser(id string, user *models.User) (*models.User, error) {
	existingUser, err := u.repos.GetUserById(id)
	if err != nil {
		return nil, err
	}
	updateTime := time.Now()
	existingUser = &models.User{
		Id:         existingUser.Id,
		Name:       user.Name,
		Surname:    user.Surname,
		SecondName: user.SecondName,
		Password:   existingUser.Password,
		Email:      existingUser.Email,
		IsVerified: user.IsVerified,
		TelegramId: user.TelegramId,
		PhotoUrl:   user.PhotoUrl,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  &updateTime,
	}

	return u.repos.UpdateUser(user)
}

func (u UserServiceImpl) DeleteUser(id string) error {
	return u.repos.DeleteUser(id)
}

func (u UserServiceImpl) ComparePassword(userPassword, inputPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(inputPassword)); err != nil {
		return errors.New("invalid email or password")
	}
	return nil
}

func (u UserServiceImpl) GenerateToken(user *models.User) (string, error) {
	cfg := config.GetConfig()

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(cfg.Jwt.Expires)).Unix(),
		Issuer:    user.Id,
	})

	return claims.SignedString([]byte(cfg.Jwt.Secret))
}

func (u UserServiceImpl) HashingPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

type UserServiceImpl struct {
	repos repository.UserRepository
}

func NewUserService(repos repository.UserRepository) UserService {
	return &UserServiceImpl{repos: repos}
}
