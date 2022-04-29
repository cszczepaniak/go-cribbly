package s3

import (
	"context"
	"io"
	"runtime"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"golang.org/x/sync/errgroup"
)

type s3Client interface {
	ListKeys(prefix string) ([]string, error)
	Download(w io.WriterAt, key string) error
	Upload(key string, body io.Reader) error
}

type rawClient struct {
	bucket string
	api    s3iface.S3API
	dl     *s3manager.Downloader
	ul     *s3manager.Uploader
}

func (c *rawClient) ListKeys(prefix string) ([]string, error) {
	out, err := c.api.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(c.bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, err
	}

	res := make([]string, 0, len(out.Contents))
	for _, obj := range out.Contents {
		res = append(res, *obj.Key)
	}
	return res, nil
}

func (c *rawClient) Download(w io.WriterAt, key string) error {
	_, err := c.dl.Download(w, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	return err
}

func (c *rawClient) Upload(key string, body io.Reader) error {
	_, err := c.ul.Upload(&s3manager.UploadInput{
		Key:         aws.String(key),
		Bucket:      aws.String(c.bucket),
		ContentType: aws.String(`application/json`),
		Body:        body,
	})
	return err
}

type ByteStore interface {
	Get(key string) ([]byte, error)
	GetWithPrefix(pref string) (map[string][]byte, error)
	Put(key string, r io.Reader) error
}

type s3ByteStore struct {
	client  s3Client
	timeout time.Duration
}

func NewS3ByteStore(bucket string, s *session.Session, timeout time.Duration) *s3ByteStore {
	ul := s3manager.NewUploader(s)
	dl := s3manager.NewDownloader(s)
	rawClient := &rawClient{
		bucket: bucket,
		dl:     dl,
		ul:     ul,
		api:    ul.S3,
	}
	return newS3Client(rawClient, timeout)
}

func newS3Client(c s3Client, timeout time.Duration) *s3ByteStore {
	return &s3ByteStore{
		client:  c,
		timeout: timeout,
	}
}

func (c *s3ByteStore) Get(key string) ([]byte, error) {
	buff := aws.NewWriteAtBuffer(nil)
	err := c.client.Download(buff, key)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func (c *s3ByteStore) GetWithPrefix(pref string) (map[string][]byte, error) {
	keys, err := c.client.ListKeys(pref)
	if err != nil {
		return nil, err
	}

	n := runtime.NumCPU()
	if n > len(keys) {
		n = len(keys)
	}

	batchSize := (len(keys) + n - 1) / n
	res := make(map[string][]byte, len(keys))
	lock := sync.Mutex{}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	for start, end := 0, batchSize; start < len(keys); start, end = end, end+batchSize {
		if end > len(keys) {
			end = len(keys)
		}
		ks := keys[start:end]
		g.Go(func() error {
			var buff []byte
			for _, k := range ks {
				select {
				case <-ctx.Done():
					return nil
				default:
				}

				w := aws.NewWriteAtBuffer(buff)
				err := c.client.Download(w, k)
				if err != nil {
					return err
				}

				lock.Lock()
				res[k] = w.Bytes()
				lock.Unlock()

				buff = buff[:]
			}
			return nil
		})
	}

	err = g.Wait()
	if err != nil {
		return nil, err
	}

	err = ctx.Err()
	if err == context.DeadlineExceeded {
		return nil, err
	}

	return res, nil
}

func (c *s3ByteStore) Put(key string, r io.Reader) error {
	return c.client.Upload(key, r)
}
