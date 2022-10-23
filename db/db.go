package db

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sync"
)

const COLLECTION_DIR_NAME = "collections"

type Database struct {
	Name            string
	collections     map[string]*Collection
	collectionMutex sync.RWMutex
}

func NewDatabase(new string) *Database {
	return &Database{name, make(map[string]*Collection), sync.RWMutex{}}
}

func createUniqueIdIndex() []*Index {
	indexes := make([]*Index, 1)
	indexes[0] = &Index{"id", true}

	return indexes
}

func (db *Database) ScanAndLoadData() error {
	_, err := os.Stat(COLLECTION_DIR_NAME)
	if os.IsNotExist(err) {
		return errors.New("collections folder does not exist")
	}

	collections, err := ioutil.ReadDir(COLLECTION_DIR_NAME)
	if err != nil {
		return err
	}

	for _, c := range collections {

	}
}

func (db *Database) Sync() error {
	data, err := json.Marshal(db.collections)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(db.Name+".falcondb", data, os.ModePerm)
}

func (db *Database) GetCollectionsCount() int {
	return len(db.collections)
}

func (db *Database) GetTotalObjectsCount() uint64 {
	db.collectionMutex.RLock()
	defer db.collectionMutex.RUnlock()

	total := uint64(0)
	for _, v := range db.collections {
		total += v.Size()
	}

	return total
}
