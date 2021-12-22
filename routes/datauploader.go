package routes

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")

	fileName := handler.Filename

	upload_filepath := filepath.Join(os.Getenv("FILE_STORAGE"), "upload")

	if _, err := os.Stat(upload_filepath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(upload_filepath, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	f, err := os.OpenFile(filepath.Join(upload_filepath, fileName), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, _ = io.WriteString(w, "File "+fileName+" Uploaded successfully")
	_, _ = io.Copy(f, file)
}
