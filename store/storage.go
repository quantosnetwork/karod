package store

import (
	"go.uber.org/zap"
)

type StorageOptions struct {
	option []*Option
	values map[string]interface{}
	libs   map[string]map[string]interface{}
}

func (so *StorageOptions) Set(key string, value string) {

	so.values[key] = value
	so.option = append(so.option)

}

type StorageOption func() string

type Option struct {
	funcs map[string]StorageOption
}

func (sof StorageOption) _get(so *StorageOptions, key string) interface{} {
	return so.values[key]
}

func (sof StorageOption) _getLibs(so *StorageOptions, key string) map[string]interface{} {
	return so.libs[key]
}

type Storage interface {
	Initialize()
	Open() (interface{}, error)
}

func (so *StorageOptions) StorageType() interface{} {
	key := "storage_type"
	option := &Option{}
	return option.funcs[key]._get(so, key)
}

func (so *StorageOptions) StorageLibs() map[string]interface{} {
	key := "storage_libs"
	option := &Option{}
	return option.funcs[key]._getLibs(so, key)
}

type StorageManager struct {
	Storage
	OpenedDB interface{}
	Logger   *zap.Logger
}

func (sm StorageManager) Initialize() {
	so := &StorageOptions{}
	libs := so.StorageLibs()
	sType := so.StorageType()
	libToLoad := libs[sType.(string)].(Storage)
	sm.OpenedDB, _ = libToLoad.Open()
}
