package common

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func SMD5(mark string, str string) string {
	u := md5.New()
	u.Write([]byte(mark + str))
	res := hex.EncodeToString(u.Sum(nil))
	return res
}

func GetRandomApnsId() string {
	randStr := strconv.Itoa(rand.Int())
	apns := SMD5(strconv.FormatInt(time.Now().Unix(), 10), randStr)
	return apns[:8] + "-" + apns[8:12] + "-" + apns[12:16] + "-" + apns[16:20] + "-" + apns[20:]
}

func GetFileLineNum() string {
	_, file, line, _ := runtime.Caller(1)
	return time.Now().Format("2006-01-02 15:04:05") + "▶ 我走到这里啦！" + file + "--line" + strconv.Itoa(line)
}

func ToSqlStr(i interface{}, seq string) string {
	var arr []string
	switch v := i.(type) {
	case string:
		arr = strings.Split(v, seq)
		break
	default:
		arr = v.([]string)
	}
	var res string
	res = "'" + arr[0] + "'"
	for i := 1; i < len(arr); i++ {
		res += "," + "'" + arr[i] + "'"
	}
	return res
}
