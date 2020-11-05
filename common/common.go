package common

import (
	"crypto/md5"
	"encoding/hex"
)

func SMD5(mark string, str string) string {
	u := md5.New()
	u.Write([]byte(mark + str))
	res := hex.EncodeToString(u.Sum(nil))
	return res
}
