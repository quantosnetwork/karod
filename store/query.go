package store

type Query interface {
	Get(key string) (value *Entry)
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
	AddItem(item *Entry)
	AddEncryptedItem(item *Entry, encKey, signature []byte)
	Write(batchInterface *interface{})
}

type Entry struct {
	Prefix []byte
	Key    []byte
	Value  []byte
}

type EncryptedEntry struct {
	Prefix    []byte
	Key       []byte
	Value     []byte
	encKey    []byte
	signature []byte
}
