package utils

import (
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
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

// SendCallbackWithHttp
// send callback return with new Queries,
// and with the original queries.
func SendCallbackWithHttp(urlStr string, newQueries []string) (string, error) {
	decodedURL, err := url.QueryUnescape(urlStr)
	if err != nil {
		return "", err
	}
	parsedURL, err := url.Parse(decodedURL)
	if err != nil {
		return "", err
	}
	originQueries := parsedURL.Query()
	for i := 0; i < len(newQueries); i += 2 {
		key := newQueries[i]
		value := newQueries[i+1]
		originQueries.Add(key, value)
	}
	parsedURL.RawQuery = originQueries.Encode()
	resp, err := http.Get(parsedURL.String())
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodyStr := string(body)
	return bodyStr, err
}
