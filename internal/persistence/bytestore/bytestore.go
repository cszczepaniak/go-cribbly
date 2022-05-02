package bytestore

import "io"

type ByteStore interface {
	Get(key string) ([]byte, error)
	GetWithPrefix(pref string) (map[string][]byte, error)
	Put(key string, r io.Reader) error
}
