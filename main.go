/**/
package main

import (
	"hoffmann/hoffmanndb"
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
		log.Errorln("Error loading .env file. Please make shure you set all system variables correctly")
	}

	log.Println("create db")
	hoffmanndb.NewHoffmannDB()

	log.Println("start building routes")
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(Handle404)

	r.HandleFunc("/", routes.Home)
	r.HandleFunc("/file", routes.UploadFile).Methods("POST")

	server_address := os.Getenv("SERVER_ADDRESS") + ":" + os.Getenv("SERVER_PORT")

	log.Println("start server in: " + server_address)

	log.Fatal(http.ListenAndServe(server_address, r))
}
