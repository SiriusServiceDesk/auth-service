package helpers

import (
	"math/rand"
	"time"
)

const (
	minId = 1000
	maxId = 9999
)

func GenerateConfirmCode() int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(maxId-minId) + minId
}
