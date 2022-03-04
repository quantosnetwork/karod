package storage

import (
	"context"
	ds "github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
	"karod/store"
	"net"
)

type StoragePlugin interface {
	store.Storage
	store.StorageOptions
	store.Query
	store.Iterator
}

type StorePlugin struct {
	conn      net.Listener
	datastore ds.Datastore
}

func (s *StorePlugin) Has(ctx context.Context, key ds.Key) (exists bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *StorePlugin) GetSize(ctx context.Context, key ds.Key) (size int, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *StorePlugin) Query(ctx context.Context, q query.Query) (query.Results, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StorePlugin) Put(ctx context.Context, key ds.Key, value []byte) error {
	//TODO implement me
	panic("implement me")
}

func (s *StorePlugin) Delete(ctx context.Context, key ds.Key) error {
	//TODO implement me
	panic("implement me")
}

func (s *StorePlugin) Sync(ctx context.Context, prefix ds.Key) error {
	//TODO implement me
	panic("implement me")
}

func (s *StorePlugin) Close() error {
	//TODO implement me
	panic("implement me")
}

func (s *StorePlugin) Get(ctx context.Context, key ds.Key) (value []byte, err error) {
	return nil, nil
}

func NewStoragePlugin() ds.Datastore {
	return &StorePlugin{}
}
