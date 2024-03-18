package object

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseObject(t *testing.T) {

	t.Run("valid", func(t *testing.T) {
		blobContent := "Hello world"
		blobStr := fmt.Sprintf("blob %d\000%s", len([]byte(blobContent)), blobContent)
		typ, content, err := parseObject(bytes.NewReader([]byte(blobStr)))

		assert.NoError(t, err)
		assert.Equal(t, "blob", typ)
		assert.Equal(t, []byte("Hello world"), content)
	})

	t.Run("invalid type", func(t *testing.T) {
		var buf bytes.Buffer
		zw := zlib.NewWriter(&buf)
		zw.Write([]byte("invalid 12345Hello world!"))
		zw.Close()

		_, _, err := parseObject(&buf)

		assert.Error(t, err)
	})

	t.Run("invalid size", func(t *testing.T) {
		var buf bytes.Buffer
		zw := zlib.NewWriter(&buf)
		zw.Write([]byte("blob invalidHello world!"))
		zw.Close()

		_, _, err := parseObject(&buf)

		assert.Error(t, err)
	})

	t.Run("size mismatch", func(t *testing.T) {
		var buf bytes.Buffer
		zw := zlib.NewWriter(&buf)
		zw.Write([]byte("blob 5Hello"))
		zw.Close()

		_, _, err := parseObject(&buf)

		assert.Error(t, err)
	})
}
func TestLoadByHashNew(t *testing.T) {

	t.Run("non-existent hash", func(t *testing.T) {
		hash := Hash{0x00, 0x01, 0x02, 0x03}

		_, _, err := LoadByHash(hash)

		assert.Error(t, err)
	})

	t.Run("io error", func(t *testing.T) {
		hash := Hash{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

		_, _, err := LoadByHash(hash)

		assert.Error(t, err)
	})

	t.Run("invalid zlib", func(t *testing.T) {
		hash := Hash{0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef}

		_, _, err := LoadByHash(hash)

		assert.Error(t, err)
	})

	t.Run("invalid object", func(t *testing.T) {
		hash := Hash{0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef}

		_, _, err := LoadByHash(hash)

		assert.Error(t, err)
	})

}

func TestLoadFileNew(t *testing.T) {
	t.Run("empty hash", func(t *testing.T) {
		hash := Hash{}

		_, _, err := LoadByHash(hash)

		assert.Error(t, err)
	})

	t.Run("corrupted zlib header", func(t *testing.T) {
		var buf bytes.Buffer
		zw := zlib.NewWriter(&buf)
		zw.Write([]byte("xinvalid"))
		zw.Close()

		_, _, err := LoadFile(&buf)

		assert.Error(t, err)
	})

	t.Run("truncated zlib stream", func(t *testing.T) {
		var buf bytes.Buffer
		zw := zlib.NewWriter(&buf)
		zw.Write([]byte("blob 5Hello"))
		zw.Close()

		_, _, err := LoadFile(&buf)

		assert.Error(t, err)
	})

	t.Run("invalid object type", func(t *testing.T) {
		var buf bytes.Buffer
		zw := zlib.NewWriter(&buf)
		zw.Write([]byte("invalid 5Hello"))
		zw.Close()

		_, _, err := LoadFile(&buf)

		assert.Error(t, err)
	})

	t.Run("negative object size", func(t *testing.T) {
		var buf bytes.Buffer
		zw := zlib.NewWriter(&buf)
		zw.Write([]byte("blob -5Hello"))
		zw.Close()

		_, _, err := LoadFile(&buf)

		assert.Error(t, err)
	})

	t.Run("io error", func(t *testing.T) {
		r := &errorReader{}

		_, _, err := LoadFile(r)

		assert.Error(t, err)
	})

	t.Run("invalid zlib", func(t *testing.T) {
		r := bytes.NewReader([]byte("invalid"))

		_, _, err := LoadFile(r)

		assert.Error(t, err)
	})

	t.Run("invalid object", func(t *testing.T) {
		var buf bytes.Buffer
		zw := zlib.NewWriter(&buf)
		zw.Write([]byte("invalid"))
		zw.Close()

		_, _, err := LoadFile(&buf)

		assert.Error(t, err)
	})

}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error")
}
