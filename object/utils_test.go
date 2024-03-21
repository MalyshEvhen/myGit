package object

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreFromFile(t *testing.T) {
	dryRun := false

	t.Run("valid file", func(t *testing.T) {
		expectedHash, err := HashFromString("a5c19667710254f835085b99726e523457150e03")
		assert.NoError(t, err)

		hash, err := StoreFromFile("test_files/test.txt", "blob", dryRun)

		assert.NoError(t, err)
		assert.Equal(t, expectedHash, hash)
	})

	t.Run("non-existing file", func(t *testing.T) {
		_, err := StoreFromFile("invalid.txt", "blob", dryRun)

		assert.Error(t, err)
	})

	t.Run("invalid type", func(t *testing.T) {
		_, err := StoreFromFile("test_files/test.txt", "invalid", dryRun)

		assert.Error(t, err)
	})
}

func TestStoreFromFileEdgeCases(t *testing.T) {
	dryRun := false

	t.Run("empty file", func(t *testing.T) {
		hash, err := StoreFromFile("", "blob", dryRun)
		assert.Error(t, err)
		assert.Equal(t, Hash{}, hash)
	})

	t.Run("large file", func(t *testing.T) {
		data := make([]byte, 1024*1024)
		err := os.WriteFile("./test_files/large.txt", data, 0644)
		require.NoError(t, err)

		hash, err := StoreFromFile("test_files/large.txt", "blob", dryRun)
		assert.NoError(t, err)
		assert.NotEqual(t, Hash{}, hash)
	})

	t.Run("invalid object type", func(t *testing.T) {
		hash, err := StoreFromFile("test_files/test.txt", "invalid", dryRun)
		assert.Error(t, err)
		assert.Equal(t, Hash{}, hash)
	})

}

func TestEncodeObjectEdgeCases(t *testing.T) {

	t.Run("empty reader", func(t *testing.T) {
		var buf bytes.Buffer
		_, err := EncodeObject(&buf, bytes.NewReader([]byte{}), "blob", 0)
		assert.NoError(t, err)
	})

	t.Run("size mismatch", func(t *testing.T) {
		var buf bytes.Buffer
		_, err := EncodeObject(&buf, bytes.NewReader([]byte("abc")), "blob", 10)
		assert.Error(t, err)
	})

	t.Run("invalid object type", func(t *testing.T) {
		var buf bytes.Buffer
		_, err := EncodeObject(&buf, bytes.NewReader([]byte("abc")), "invalid", 3)
		assert.Error(t, err)
	})

}

func TestCompressObjectEdgeCases(t *testing.T) {
	// Test compressing an e valid input
	t.Run("valid input", func(t *testing.T) {
		var buf bytes.Buffer
		r := []byte("Hello world")

		err := Compress(&buf, r)
		assert.NoError(t, err)
		assert.NotEqual(t, r, buf.Bytes())
		require.Greater(t, len(buf.Bytes()), 0)

		// Test decompression
		bytesToDecompress := bytes.NewReader(buf.Bytes())
		zr, err := zlib.NewReader(bytesToDecompress)
		assert.NoError(t, err)
		decompressed := make([]byte, len(r))
		_, err = zr.Read(decompressed)
		fmt.Println("decompressed:", string(decompressed))
		assert.NoError(t, err)
		assert.Equal(t, r, decompressed)
	})

	// Test compressing an empty input
	t.Run("empty input", func(t *testing.T) {
		var buf bytes.Buffer
		err := Compress(&buf, []byte{})
		assert.NoError(t, err)
		assert.Equal(t, 0, len(buf.Bytes()))
	})

	// Test compressing and decompressing a valid input
	t.Run("compress and decompress", func(t *testing.T) {

		tests := []struct {
			name  string
			input []byte
			want  []byte
		}{
			{"hello world", []byte("hello world"), []byte("hello world")},
			{"binary data", []byte{0x1, 0x2, 0x3}, []byte{0x1, 0x2, 0x3}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var buf bytes.Buffer
				err := Compress(&buf, tt.input)
				assert.NoError(t, err)

				decompressed := make([]byte, len(tt.input))
				reader := bytes.NewReader(buf.Bytes())
				n, err := zlib.NewReader(reader)
				assert.NoError(t, err)

				i, err := n.Read(decompressed)
				assert.NoError(t, err)

				assert.Equal(t, tt.want, decompressed[:i])
				buf.Reset()
			})
		}
	})
}

func BenchmarkCompress(b *testing.B) {
	data := make([]byte, 1024)
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		err := Compress(&buf, data)
		require.NoError(b, err)
	}
}

func FuzzCompress(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var buf bytes.Buffer
		err := Compress(&buf, data)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func FuzzDecompress(f *testing.F) {
	f.Fuzz(func(t *testing.T, compressed []byte) {
		decompressed := make([]byte, 100)
		reader := bytes.NewReader(compressed)
		n, err := zlib.NewReader(reader)
		if err != nil {
			return
		}

		_, err = n.Read(decompressed)
		if err != nil {
			t.Fatal(err)
		}
	})
}
