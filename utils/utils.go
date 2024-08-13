package utils

import (
	"os"
	"time"
)

func FileIsExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return err == nil
}

func IsNotEmptyStr(s string) bool {
	if s != "" {
		return true
	} else {
		return false
	}
}

func NowCST() string {
	location, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(location).Format(time.DateTime)
}
