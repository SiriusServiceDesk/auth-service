package helpers

import (
	"math/rand"
	"time"
)

const (
	minId = 100000
	maxId = 999999
)

func GenerateConfirmCode() int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(maxId-minId) + minId
}
