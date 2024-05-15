package helpers

import (
	"errors"
	"github.com/SiriusServiceDesk/auth-service/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"strings"
)

func getClaims(tokenString string) (*jwt.StandardClaims, error) {
	cfg := config.GetConfig()
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Jwt.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims := token.Claims.(*jwt.StandardClaims)

	return claims, err
}

func GetUserIdFromToken(ctx *fiber.Ctx) (string, error) {
	authHeader := ctx.GetReqHeaders()[fiber.HeaderAuthorization]

	if len(authHeader) == 0 {
		return "", errors.New("auth token is required")
	}

	if authHeader[0] == "" {
		return "", errors.New("auth token is required")
	}

	tokenSplit := strings.Split(authHeader[0], " ")

	if tokenSplit[0] != "Bearer" {
		return "", errors.New("invalid token header")
	}

	if tokenSplit[1] == "" {
		return "", errors.New("token must not be empty")
	}

	claims, err := getClaims(tokenSplit[1])
	if err != nil {
		return "", err
	}

	return claims.Issuer, nil
}
