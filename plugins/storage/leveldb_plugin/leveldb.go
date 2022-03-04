package plugins

import (
	"github.com/syndtr/goleveldb/leveldb"
	"karod/store"
	"strings"
)

/*
	LevelDB Quantos Karod Storage plugin
*/

type LevelDB struct {
	store.Storage
	store.Query
	store.Iterator
	db *leveldb.DB
}

func (ldb *LevelDB) Initialize() {
	var err error
	db, err := ldb.Open()
	if err != nil {
		panic(err)
	}
	DB = db.(*LevelDB)
}

func (ldb *LevelDB) Open() (interface{}, error) {
	db, err := leveldb.OpenFile("./data/db", nil)
	ldb.db = db
	defer ldb.db.Close()
	return ldb, err
}

func NewLevelDB() store.Storage {
	ldb := &LevelDB{}
	return ldb
}

// DB contains the LevelDB pointer
var DB *LevelDB

func init() {
	DB.Initialize()
	defer DB.db.Close()
}

func (ldb *LevelDB) Put(key string, value string) error {
	err := DB.db.Put([]byte(key), []byte(value), nil)
	if err != nil {
		return err
	}
	return nil
}

func (ldb *LevelDB) Get(key string) (value *store.Entry) {
	entry, err := DB.db.Get([]byte(key), nil)
	if err != nil {
		return &store.Entry{}
	}
	return &store.Entry{
		Key:   []byte(key),
		Value: entry,
	}
}

func (ldb *LevelDB) PutWithPrefix(prefix, key string, value string) error {
	var k []string
	k[0] = prefix
	k[1] = key
	kk := strings.Join(k, "_")
	return ldb.Put(kk, value)
}

func (ldb *LevelDB) GetWithPrefix(prefix, key string) (value *store.Entry) {
	var k []string
	k[0] = prefix
	k[1] = key
	kk := strings.Join(k, "_")
	return ldb.Get(kk)
}
