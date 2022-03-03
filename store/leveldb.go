package store

import (
	"github.com/syndtr/goleveldb/leveldb"
)

/*
	LevelDB Quantos Karod Storage plugin
*/

type LevelDB struct {
	Storage
	db *leveldb.DB
}

func (ldb *LevelDB) Open() (*LevelDB, error) {
	db, err := leveldb.OpenFile("./data/db", nil)
	ldb.db = db
	return ldb, err
}
