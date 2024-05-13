package dependencies

import (
	"github.com/SiriusServiceDesk/auth-service/internal/repository"
	"github.com/SiriusServiceDesk/auth-service/internal/services"
)

type Container struct {
	UserService services.UserService
	Redis       repository.RedisRepository
}
