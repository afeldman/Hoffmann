package routes

import (
	"hoffmann/dataloader"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
)

func HTTPUploader(w http.ResponseWriter, r *http.Request) {
	dataloader.EnsureDir(dataloader.FileStorage.Path)
	dataloader.EnsureDir(dataloader.FileStorage.TempPath)

	if r.Method != "POST" || r.Header.Get("Session-ID") == "" || r.Header.Get("Content-Range") == "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Invalid request."))
	}

	var upload dataloader.UploadFile

	sessionID := r.Header.Get("Session-ID")
	contentRange := r.Header.Get("Content-Range")

	body, err := ioutil.ReadAll(r.Body)
	dataloader.CheckError(err)

	totalSize, partFrom, partTo := dataloader.ParseContentRange(contentRange)

	if partFrom == 0 {
		_, ok := dataloader.Files[sessionID]
		if !ok {
			w.WriteHeader(http.StatusCreated)

			_, params, err := mime.ParseMediaType(r.Header.Get("Content-Disposition"))
			dataloader.CheckError(err)
			fileName := params["filename"]

			newFile := dataloader.FileStorage.TempPath + "/" + sessionID
			_, err = os.Create(newFile)
			dataloader.CheckError(err)

			f, err := os.OpenFile(newFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			dataloader.CheckError(err)

			dataloader.Files[sessionID] = dataloader.UploadFile{
				File:     f,
				Name:     fileName,
				TempPath: newFile,
				Status:   dataloader.CREATED,
				Size:     totalSize,
			}
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}

	upload = dataloader.Files[sessionID]
	upload.Status = dataloader.UPLOADING

	_, err = upload.File.Write(body)
	dataloader.CheckError(err)

	upload.File.Sync()
	upload.Transfered = partTo

	w.Header().Set("Content-Length", string(len(body)))
	w.Header().Set("Connection", "close")
	w.Header().Set("Range", contentRange)
	w.Write([]byte(contentRange))

	if partTo >= totalSize {
		dataloader.MoveToPath(sessionID)
		upload.File.Close()
		delete(dataloader.Files, sessionID)
	}

}
