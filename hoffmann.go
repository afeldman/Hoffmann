package main

import (
	"sync"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/codec/msgpack"

	log "github.com/sirupsen/logrus"
)

type database struct {
	db storm.DB
}

type HoffmannEntry struct {
	Name       string   `storm:"id"`
	Karel_File []string `storm:"karel"`
}

var (
	lock     = &sync.Mutex{}
	instance *database
)

func NewHoffmannDB() *database {

	lock.Lock()
	defer lock.Unlock()

	if instance == nil {
		// ainda não é a melhor implementação devido
		// os bloqueios
		db_, err := storm.Open("hoffmann.db", storm.Codec(msgpack.Codec))
		if err != nil {
			log.Fatal(err)
		}
		defer db_.Close()

		instance = &database{
			db: *db_,
		}
	}

	return instance

}
