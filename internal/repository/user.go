package repository

import (
	"fmt"
	"github.com/SiriusServiceDesk/auth-service/internal/config"
	"github.com/SiriusServiceDesk/auth-service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUsers() ([]*models.User, error)
	GetUser(id string) (*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUser(id string) error
}

func (u UserRepositoryImpl) GetUsers() ([]*models.User, error) {
	var result []*models.User
	if err := u.db.Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (u UserRepositoryImpl) GetUser(id string) (*models.User, error) {
	var result *models.User
	if err := u.db.Where(models.User{Id: id}).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (u UserRepositoryImpl) CreateUser(user *models.User) (*models.User, error) {
	if err := u.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u UserRepositoryImpl) UpdateUser(user *models.User) (*models.User, error) {
	if err := u.db.Save(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u UserRepositoryImpl) DeleteUser(id string) error {
	if err := u.db.Delete(models.User{Id: id}).Error; err != nil {
		return err
	}
	return nil
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	cfg := config.GetConfig()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Db.Host, cfg.Db.User, cfg.Db.Password, cfg.Db.Name, cfg.Db.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	pgSvc := &UserRepositoryImpl{db: db}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}
	return pgSvc
}
