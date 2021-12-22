package routes

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

/*
curl -X POST -i -H "Accept: * /*" \
	-H "Accept-Encoding: gzip, deflate, br" \
	-H "Connection: keep-alive" \
	-H "Content-Type: multipart/form-data" \
	-F file="<data>" http://localhost:2611/file
*/
func UploadFile(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")

	fileName := handler.Filename

	fileExtension := filepath.Ext(fileName)
	if !strings.EqualFold(fileExtension, ".karel") {
		_, _ = io.WriteString(w, "file extension "+fileExtension+" not karel type")
		return
	}

	upload_filepath := filepath.Join(os.Getenv("FILE_STORAGE"), "uploads")

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

	// f is copied now work on the karel system.
	// 1. check if kpc is included
	// 		use nishimura to check
	

}
