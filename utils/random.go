package utils

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

var seederPhone = []string{
	"1(234)5678901x1234",
	"(+351) 282 43 50 50",
	"90191919908",
	"555-8909",
	"001 6867684",
	"001 6867684x1",
	"1 (234) 567-8901",
	"1-234-567-8901 ext1234",
	"+62811132431",
	"+62-811-132-431",
	"62(751)142345",
	"6285274507699",
	"089899992834",
	"+6285274507699",
	"85274507699",
	"8527450769999",
}

func RandomSender() string {
	get := RandomInt(0, int64(len(seederPhone)-1))
	return seederPhone[get]
}

func RandomReceiver() string {
	get := RandomInt(0, int64(len(seederPhone)-1))
	return seederPhone[get]
}

func RandomBody() string {
	return RandomString(100)
}
