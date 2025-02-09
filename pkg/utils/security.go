package utils

import (
	"crypto/rand"
	"math/big"
)

// GenerateAlNumOTP - Function to generate an n-character alphanumeric OTP
// Working Sample - https://go.dev/play/p/5FUKG0YZNAP
func (u *Utils) GenerateAlNumOTP(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	otp := make([]byte, length)
	for i := range otp {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		otp[i] = charset[randomIndex.Int64()]
	}
	return string(otp), nil
}
