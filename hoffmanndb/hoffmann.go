package hoffmanndb

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/codec/msgpack"

	log "github.com/sirupsen/logrus"
)

type database struct {
	db storm.DB
}

type HoffmannEntry struct {
	Hash      string `storm:"id"`
	Name      string `storm:"name"`
	Version   string `storm:"index"`
	KarelFile string `storm:"karel"`
}

var (
	lock     = &sync.Mutex{}
	instance *database
)

func NewHoffmannDB() *database {

	lock.Lock()
	defer lock.Unlock()

	if instance == nil {
		db_, err := storm.Open(filepath.Join(os.Getenv("DATABASE_STORAGE"), "hoffmann.db"), storm.Codec(msgpack.Codec))
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

func (db *database) Save(hoffmannentry *HoffmannEntry) error {
	lock.Lock()
	defer lock.Unlock()

	return db.db.Save(hoffmannentry)
}

func (db *database) Delete(hoffmannentry *HoffmannEntry) error {
	lock.Lock()
	defer lock.Unlock()

	return db.db.DeleteStruct(hoffmannentry)
}
