package dataloader

import (
	"os"
	"time"
)

const (
	CREATED   = 0
	UPLOADING = 1
)

type STATUS int32

type UploadFile struct {
	File       *os.File
	Name       string
	TempPath   string
	Status     STATUS
	Size       int64
	Transfered int64
}

type fileStorage struct {
	Path     string
	TempPath string
}

var FileStorage = fileStorage{
	Path:     "./files",
	TempPath: ".tmp",
}

var Files = make(map[string]UploadFile)

func MoveToPath(id string) {
	uploadFile := Files[id]
	filePath := FileStorage.Path + "/" + uploadFile.Name
	if FileExists(filePath) {
		t := time.Now().Format(time.RFC3339)
		filePath = FileStorage.Path + "/" + t + "-" + uploadFile.Name
	}

	err := os.Rename(uploadFile.TempPath, filePath)
	CheckError(err)
}
