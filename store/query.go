package store

type Query interface {
	Get(key string) (value *interface{})
	Put(key string, value string) error
	Delete(key string) error
}

type Iterator interface {
	New() *interface{}
	Next()
	Seek(key string)
	Release()
}

type Batch interface {
	NewBatch() *interface{}
	Write(batchInterface *interface{})
}
