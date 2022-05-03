package bytestore

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cszczepaniak/go-cribbly/internal/random"
)

func randomReader(t *testing.T) ([]byte, io.Reader) {
	u, err := uuid.NewUUID()
	require.NoError(t, err)

	return u[:], bytes.NewReader(u[:])
}

func TestGetTimeout(t *testing.T) {
	rawClient := newMemoryRawClient()
	client := newS3Client(rawClient, time.Millisecond)

	k := random.UUID()
	_, r := randomReader(t)
	client.Put(k, r)

	bs, err := client.Get(k)
	require.NoError(t, err)
	require.NotEmpty(t, bs)

	rawClient.addDelay(time.Second)

	bs, err = client.Get(k)
	require.Equal(t, context.DeadlineExceeded, err)
	require.Empty(t, bs)
}

func TestPutTimeout(t *testing.T) {
	rawClient := newMemoryRawClient()
	client := newS3Client(rawClient, time.Millisecond)

	rawClient.addDelay(time.Second)

	k := random.UUID()
	_, r := randomReader(t)

	err := client.Put(k, r)
	require.Equal(t, context.DeadlineExceeded, err)
}

func TestGetWithPrefix(t *testing.T) {
	rawClient := newMemoryRawClient()
	client := newS3Client(rawClient, time.Second)

	t.Run(`zero keys`, func(t *testing.T) {
		res, err := client.GetWithPrefix(`a`)
		require.NoError(t, err)
		assert.Empty(t, res)
	})

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

	t.Run(`timeout`, func(t *testing.T) {
		t.Cleanup(rawClient.clear)

		client = newS3Client(rawClient, time.Millisecond)

		for i := 0; i < 5; i++ {
			_, r := randomReader(t)
			rawClient.Upload(fmt.Sprintf(`a%d.json`, i), r)
		}

		rawClient.addDelay(20 * time.Millisecond)

		res, err := client.GetWithPrefix(`a`)
		assert.Equal(t, context.DeadlineExceeded, err)
		assert.Empty(t, res)
	})
}

type memoryClient struct {
	lock       sync.Mutex
	objects    map[string][]byte
	errorOnKey string
	delay      time.Duration
}

var _ s3Client = (*memoryClient)(nil)

func newMemoryRawClient() *memoryClient {
	return &memoryClient{
		objects: make(map[string][]byte),
	}
}

func (c *memoryClient) clear() {
	c.errorOnKey = ``
	c.delay = 0
	c.lock.Lock()
	defer c.lock.Unlock()
	c.objects = make(map[string][]byte)
}

func (c *memoryClient) setErrorKey(k string) {
	c.errorOnKey = k
}

func (c *memoryClient) addDelay(d time.Duration) {
	c.delay = d
}

func (c *memoryClient) ListKeys(prefix string) ([]string, error) {
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

func (c *memoryClient) Download(w io.WriterAt, key string) error {
	time.Sleep(c.delay)
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

func (c *memoryClient) DownloadWithContext(ctx context.Context, w io.WriterAt, key string) error {
	done := make(chan struct{})
	var err error
	go func() {
		err = c.Download(w, key)
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return err
	}
}

func (c *memoryClient) Upload(key string, body io.Reader) error {
	time.Sleep(c.delay)
	bs, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	c.objects[key] = bs

	return nil
}

func (c *memoryClient) UploadWithContext(ctx context.Context, key string, body io.Reader) error {
	done := make(chan struct{})
	var err error
	go func() {
		err = c.Upload(key, body)
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return err
	}
}
