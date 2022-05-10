package bytestore

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"sync"

	"github.com/cszczepaniak/go-cribbly/internal/cribblyerr"
)

type MemoryByteStore struct {
	blobs map[string][]byte
	lock  sync.Mutex
}

var _ ByteStore = (*MemoryByteStore)(nil)

func NewMemoryByteStore() *MemoryByteStore {
	return &MemoryByteStore{
		blobs: make(map[string][]byte),
	}
}

func (m *MemoryByteStore) Get(key string) ([]byte, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	blob, ok := m.blobs[key]
	if !ok {
		return nil, cribblyerr.ErrNotFound(key)
	}
	return blob, nil
}

func (m *MemoryByteStore) GetJSON(key string, v interface{}) error {
	bs, err := m.Get(key)
	if err != nil {
		return err
	}

	return json.Unmarshal(bs, v)
}

func (m *MemoryByteStore) GetWithPrefix(pref string) (map[string][]byte, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	res := make(map[string][]byte)
	for k, blob := range m.blobs {
		if strings.HasPrefix(k, pref) {
			res[k] = blob
		}
	}

	return res, nil
}

func (m *MemoryByteStore) Put(key string, r io.Reader) error {
	bs, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	m.blobs[key] = bs
	return nil
}

func (m *MemoryByteStore) PutJSON(key string, v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return m.Put(key, bytes.NewReader(bs))
}
