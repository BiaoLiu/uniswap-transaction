package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())

}

func Randn(n int) int {
	return rand.Intn(n)
}
