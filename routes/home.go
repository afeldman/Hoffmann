package routes

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func Home(w http.ResponseWriter, r *http.Request) {

	log.Println("nothing new only home :D")
	w.Write([]byte("home"))
}
