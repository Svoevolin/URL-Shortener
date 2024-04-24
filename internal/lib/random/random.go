package random

import (
	"time"

	"golang.org/x/exp/rand"
)

// NewRandomStrings generates random string with given size
func NewRandomStrings(size int8) string {
	rnd := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")

	b := make([]rune, size)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b)
}
