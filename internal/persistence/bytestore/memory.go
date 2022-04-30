// +build !prod

package bytestore

import (
	"errors"
	"io"
	"strings"
	"sync"
	"time"
)

type MemoryS3Client struct {
	lock       sync.Mutex
	objects    map[string][]byte
	errorOnKey string
	delay      time.Duration
}

var _ S3Client = (*MemoryS3Client)(nil)

func NewMemoryS3Client() *MemoryS3Client {
	return &MemoryS3Client{
		objects: make(map[string][]byte),
	}
}

func (c *MemoryS3Client) Clear() {
	c.errorOnKey = ``
	c.delay = 0
	c.lock.Lock()
	defer c.lock.Unlock()
	c.objects = make(map[string][]byte)
}

func (c *MemoryS3Client) SetErrorKey(k string) {
	c.errorOnKey = k
}

func (c *MemoryS3Client) AddDelay(d time.Duration) {
	c.delay = d
}

func (c *MemoryS3Client) ListKeys(prefix string) ([]string, error) {
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

func (c *MemoryS3Client) Download(w io.WriterAt, key string) error {
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

func (c *MemoryS3Client) Upload(key string, body io.Reader) error {
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
