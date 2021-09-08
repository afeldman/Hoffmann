package dataloader

import (
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func EnsureDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
}

func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		return true
	}

	return false
}

func generateContentRange(index uint64, fileChunk int, partSize int, totalSize int64) string {
	var contentRange string
	if index == 0 {
		contentRange = "bytes 0-" + fmt.Sprintf("%v", partSize) + "/" + fmt.Sprintf("%v", totalSize)
	} else {
		from := uint64(fileChunk) * index
		to := uint64(fileChunk) * (index + 1)
		contentRange = "bytes " + fmt.Sprintf("%v", from) + "-" + fmt.Sprintf("%v", to) + "/" + fmt.Sprintf("%v", totalSize)
	}

	return contentRange
}

func parseBody(body string) int64 {
	fromTo := strings.Split(body, "/")[0]
	splitted := strings.Split(fromTo, "-")

	partTo, err := strconv.ParseInt(splitted[1], 10, 64)
	CheckError(err)

	return partTo
}

func generateSessionID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%X", b)
}

func ParseContentRange(contentRange string) (totalSize int64, partFrom int64, partTo int64) {
	contentRange = strings.Replace(contentRange, "bytes ", "", -1)
	fromTo := strings.Split(contentRange, "/")[0]
	totalSize, err := strconv.ParseInt(strings.Split(contentRange, "/")[1], 10, 64)
	CheckError(err)

	splitted := strings.Split(fromTo, "-")

	partFrom, err = strconv.ParseInt(splitted[0], 10, 64)
	CheckError(err)
	partTo, err = strconv.ParseInt(splitted[1], 10, 64)
	CheckError(err)

	return totalSize, partFrom, partTo
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
