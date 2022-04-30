package bytestore

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"runtime"
	"testing"
	"time"

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
	rawClient := NewMemoryS3Client()
	client := newS3Client(rawClient, time.Second)

	t.Run(`zero keys`, func(t *testing.T) {
		res, err := client.GetWithPrefix(`a`)
		require.NoError(t, err)
		assert.Empty(t, res)
	})

	t.Run(`a few`, func(t *testing.T) {
		t.Cleanup(rawClient.Clear)
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
		t.Cleanup(rawClient.Clear)
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
		t.Cleanup(rawClient.Clear)
		n := 3*runtime.NumCPU() + 7
		for i := 0; i < n; i++ {
			_, r := randomReader(t)
			rawClient.Upload(fmt.Sprintf(`a%d.json`, i), r)
		}

		errKey := fmt.Sprintf(`a%d.json`, rand.Int31n(int32(n)))
		rawClient.SetErrorKey(errKey)

		res, err := client.GetWithPrefix(`a`)
		assert.Equal(t, errors.New(`configured error`), err)
		assert.Empty(t, res)
	})

	t.Run(`timeout`, func(t *testing.T) {
		t.Cleanup(rawClient.Clear)

		client = newS3Client(rawClient, time.Millisecond)

		for i := 0; i < 5; i++ {
			_, r := randomReader(t)
			rawClient.Upload(fmt.Sprintf(`a%d.json`, i), r)
		}

		rawClient.AddDelay(20 * time.Millisecond)

		res, err := client.GetWithPrefix(`a`)
		assert.Equal(t, context.DeadlineExceeded, err)
		assert.Empty(t, res)
	})
}
