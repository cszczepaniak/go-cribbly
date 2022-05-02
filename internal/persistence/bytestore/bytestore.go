package bytestore

import "io"

type ByteStore interface {
	Get(key string) ([]byte, error)
	GetJSON(key string, v interface{}) error
	Put(key string, r io.Reader) error
	PutJSON(key string, v interface{}) error
	GetWithPrefix(pref string) (map[string][]byte, error)
}
