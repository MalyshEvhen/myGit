package object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashFromString(t *testing.T) {
	t.Run("valid hash", func(t *testing.T) {
		hashStr := "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"
		expectedHash := Hash{0xa9, 0x4a, 0x8f, 0xe5, 0xcc, 0xb1, 0x9b, 0xa6, 0x1c, 0x4c, 0x08, 0x73, 0xd3, 0x91, 0xe9, 0x87, 0x98, 0x2f, 0xbb, 0xd3}

		hash, err := HashFromString(hashStr)
		assert.NoError(t, err)
		assert.Equal(t, expectedHash, hash)
	})

	t.Run("invalid length", func(t *testing.T) {
		hashStr := "a94a8fe5ccb19ba61c4c0873d391e987982fbbd" // missing last 2 chars

		_, err := HashFromString(hashStr)
		assert.Error(t, err)
	})

	t.Run("invalid hex", func(t *testing.T) {
		hashStr := "a94a8fe5ccb19ba61c4c0873d391e987982fbxg" // invalid hex digits

		_, err := HashFromString(hashStr)
		assert.Error(t, err)
	})
}
func TestHashFromStringNew(t *testing.T) {

	t.Run("empty string", func(t *testing.T) {
		hashStr := ""

		_, err := HashFromString(hashStr)

		assert.Error(t, err)
	})

	t.Run("all zero hash", func(t *testing.T) {
		hashStr := "0000000000000000000000000000000000000000"

		expectedHash := Hash{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

		hash, err := HashFromString(hashStr)

		assert.NoError(t, err)
		assert.Equal(t, expectedHash, hash)
	})

	t.Run("max value hash", func(t *testing.T) {
		hashStr := "ffffffffffffffffffffffffffffffffffffffff"

		expectedHash := Hash{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

		hash, err := HashFromString(hashStr)

		assert.NoError(t, err)
		assert.Equal(t, expectedHash, hash)
	})

	t.Run("invalid hex character", func(t *testing.T) {
		hashStr := "a94a8fe5ccb19ba61c4c0873d391e987982fbxgd3"

		_, err := HashFromString(hashStr)

		assert.Error(t, err)
	})

}

func TestHashString(t *testing.T) {

	t.Run("empty hash", func(t *testing.T) {
		var hash Hash

		hashStr := hash.String()

		assert.Equal(t, "0000000000000000000000000000000000000000", hashStr)
	})

	t.Run("max value hash", func(t *testing.T) {
		hash := Hash{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

		expected := "ffffffffffffffffffffffffffffffffffffffff"

		hashStr := hash.String()

		assert.Equal(t, expected, hashStr)
	})

}
