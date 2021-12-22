package queue

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/afeldman/kpc"
	log "github.com/sirupsen/logrus"

	"github.com/chigopher/pathlib"
	archiver "github.com/mholt/archiver/v3"
	"github.com/otium/queue"
)

var (
	instance *queue.Queue
	doOnce   sync.Once
)

func HoffmannQueue() *queue.Queue {

	doOnce.Do(func() {

		instance = queue.NewQueue(func(val interface{}) {
			//build temp dir
			dir, err := ioutil.TempDir("hoffmann", "archive_")
			if err != nil {
				log.Fatal(err)
			}
			defer os.RemoveAll(dir)

			dir_path := pathlib.NewPath(dir)

			// filepath
			filepath := pathlib.NewPath(val.(string))

			// if the file exists or another error happends stop
			exists, err := filepath.Exists()
			if err != nil || !exists {
				return
			}

			// unpack data
			err = archiver.Unarchive(filepath.String(), dir)
			if err != nil {
				return
			}

			// if it is not kpc in the root, then delete the file
			kpc_filepath := dir_path.Join(filepath.Name() + ".kpc")
			exists, err = kpc_filepath.Exists()
			if err != nil || !exists {
				filepath.RemoveAll()
				return
			}

			// parse KPC
			kpc_err, kpc_info := kpc.ReadKPCFile(kpc_filepath.String())
			if kpc_err != nil || !exists {
				filepath.RemoveAll()
				return
			}

			fmt.Println(kpc_info.Name)

		}, 100) // hundred strings are allowed.
	})

	return instance
}
