package common

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(s string) string {
	x := md5.Sum([]byte(s))
	return hex.EncodeToString(x[:])
}