package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Encrypt(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}
