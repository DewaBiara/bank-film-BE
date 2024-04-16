package impl

import (
	"math/rand"
)

func GenerateRandomOTP() string {
	var otpsets = []rune("0123456789")
	otps := make([]rune, 6)
	for i := range otps {
		otps[i] = otpsets[rand.Intn(len(otpsets))]
	}
	return string(otps)
}
