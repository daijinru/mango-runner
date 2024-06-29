package utils

import (
	"github.com/google/uuid"
	"time"
)

func TimeNow() string {
	return time.Now().Format("01-02-2006 15:04:05")
}

func GenerateUUIDFileName() string {
	u := uuid.New()

	fileName := u.String()
	fileName = fileName[:8] + fileName[9:13] + fileName[14:18] + fileName[19:23] + fileName[24:]

	return fileName
}

// ConvertArrayToStr convert []string, output 1 string
func ConvertArrayToStr(arr []string) string {
	merged := ""
	for i := 0; i < len(arr); i++ {
		merged += arr[i]
	}
	return merged
}
