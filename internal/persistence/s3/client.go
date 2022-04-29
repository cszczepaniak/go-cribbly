package s3

import (
	"errors"
	"io"
	"runtime"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type rawS3Client interface {
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

type s3Client struct {
	rawClient rawS3Client
}

func NewS3Client(bucket string, s *session.Session) *s3Client {
	ul := s3manager.NewUploader(s)
	dl := s3manager.NewDownloader(s)
	rawClient := &rawClient{
		bucket: bucket,
		dl:     dl,
		ul:     ul,
		api:    ul.S3,
	}
	return newS3Client(rawClient)
}

func newS3Client(c rawS3Client) *s3Client {
	return &s3Client{
		rawClient: c,
	}
}

func (c *s3Client) Get(key string) ([]byte, error) {
	buff := aws.NewWriteAtBuffer(nil)
	err := c.rawClient.Download(buff, key)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func (c *s3Client) GetWithPrefix(pref string) (map[string][]byte, error) {
	keys, err := c.rawClient.ListKeys(pref)
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
	cancel := make(chan struct{})

	wg := sync.WaitGroup{}
	errs := make(chan error)
	for start, end := 0, batchSize; start < len(keys); start, end = end, end+batchSize {
		wg.Add(1)
		pd := &parallelDownloader{
			client: c.rawClient,
			cancel: make(<-chan struct{}),
			errs:   errs,
		}
		if end > len(keys) {
			end = len(keys)
		}
		go pd.doDownload(&wg, keys[start:end], func(k string, bs []byte) {
			lock.Lock()
			defer lock.Unlock()
			res[k] = bs
		})
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case err := <-errs:
		close(cancel)
		return nil, err
	case <-done:
		return res, nil
	case <-time.After(time.Minute):
		close(cancel)
		return nil, errors.New(`request timed out`)
	}
}

type parallelDownloader struct {
	client rawS3Client
	cancel <-chan struct{}
	errs   chan<- error
}

func (d *parallelDownloader) doDownload(wg *sync.WaitGroup, keys []string, onComplete func(k string, bs []byte)) {
	defer wg.Done()
	var buff []byte
	for _, k := range keys {
		select {
		case <-d.cancel:
			return
		default:
		}

		w := aws.NewWriteAtBuffer(buff)
		err := d.client.Download(w, k)
		if err != nil {
			d.errs <- err
			return
		}
		onComplete(k, w.Bytes())
		buff = buff[:]
	}
}

func (c *s3Client) Put(key string, r io.Reader) error {
	return c.rawClient.Upload(key, r)
}
