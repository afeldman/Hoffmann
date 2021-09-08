package main

import (
	"hoffmann/routes"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func Handle404(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("custom 404"))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(Handle404)
	r.HandleFunc("/", routes.Home)
	r.HandleFunc("/data", routes.HTTPUploader)

	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_ADDRESS")+":"+os.Getenv("SERVER_PORT"), r))
}
