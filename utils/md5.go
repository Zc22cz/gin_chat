package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tempstr := h.Sum(nil)
	return hex.EncodeToString(tempstr)
}

func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

func MakePassword(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)
}

func ValidPassword(plainpwd, salt, password string) bool {
	return Md5Encode(plainpwd+salt) == password
}
