package s3

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func randomReader(t *testing.T) ([]byte, io.Reader) {
	u, err := uuid.NewUUID()
	require.NoError(t, err)

	return u[:], bytes.NewReader(u[:])
}

func TestGetWithPrefix(t *testing.T) {
	rawClient := newMemoryRawClient()
	client := newS3Client(rawClient)

	t.Run(`a few`, func(t *testing.T) {
		t.Cleanup(rawClient.clear)
		a1, a1R := randomReader(t)
		rawClient.Upload(`a1.json`, a1R)
		a2, a2R := randomReader(t)
		rawClient.Upload(`a2.json`, a2R)
		a3, a3R := randomReader(t)
		rawClient.Upload(`a3.json`, a3R)
		_, b1R := randomReader(t)
		rawClient.Upload(`b1.json`, b1R)

		res, err := client.GetWithPrefix(`a`)
		require.NoError(t, err)
		assert.Len(t, res, 3)
		assert.Contains(t, res, `a1.json`)
		assert.Equal(t, a1, res[`a1.json`])
		assert.Contains(t, res, `a2.json`)
		assert.Equal(t, a2, res[`a2.json`])
		assert.Contains(t, res, `a3.json`)
		assert.Equal(t, a3, res[`a3.json`])
	})

	t.Run(`more than num CPU`, func(t *testing.T) {
		t.Cleanup(rawClient.clear)
		n := 3*runtime.NumCPU() + 7
		for i := 0; i < n; i++ {
			_, r := randomReader(t)
			rawClient.Upload(fmt.Sprintf(`a%d.json`, i), r)
		}

		res, err := client.GetWithPrefix(`a`)
		require.NoError(t, err)
		assert.Len(t, res, n)
		for i := 0; i < n; i++ {
			assert.Contains(t, res, fmt.Sprintf(`a%d.json`, i))
		}
	})

	t.Run(`error partway through`, func(t *testing.T) {
		t.Cleanup(rawClient.clear)
		n := 3*runtime.NumCPU() + 7
		for i := 0; i < n; i++ {
			_, r := randomReader(t)
			rawClient.Upload(fmt.Sprintf(`a%d.json`, i), r)
		}

		errKey := fmt.Sprintf(`a%d.json`, rand.Int31n(int32(n)))
		rawClient.setErrorKey(errKey)

		res, err := client.GetWithPrefix(`a`)
		assert.Equal(t, errors.New(`configured error`), err)
		assert.Empty(t, res)
	})
}

type memoryRawClient struct {
	lock       sync.Mutex
	objects    map[string][]byte
	errorOnKey string
}

var _ rawS3Client = (*memoryRawClient)(nil)

func newMemoryRawClient() *memoryRawClient {
	return &memoryRawClient{
		objects: make(map[string][]byte),
	}
}

func (c *memoryRawClient) clear() {
	c.errorOnKey = ``
	c.lock.Lock()
	defer c.lock.Unlock()
	c.objects = make(map[string][]byte)
}

func (c *memoryRawClient) setErrorKey(k string) {
	c.errorOnKey = k
}

func (c *memoryRawClient) ListKeys(prefix string) ([]string, error) {
	var res []string

	c.lock.Lock()
	defer c.lock.Unlock()
	for k := range c.objects {
		if strings.HasPrefix(k, prefix) {
			res = append(res, k)
		}
	}
	return res, nil
}

func (c *memoryRawClient) Download(w io.WriterAt, key string) error {
	if c.errorOnKey != `` && key == c.errorOnKey {
		return errors.New(`configured error`)
	}
	c.lock.Lock()
	defer c.lock.Unlock()

	v, ok := c.objects[key]
	if !ok {
		return errors.New(`not found`)
	}

	_, err := w.WriteAt(v, 0)
	return err
}

func (c *memoryRawClient) Upload(key string, body io.Reader) error {
	bs, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	c.objects[key] = bs

	return nil
}
