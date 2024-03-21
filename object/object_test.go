package object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewObject(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		o, err := NewObject("blob", 10, []byte("hello"))
		assert.NoError(t, err)
		assert.Equal(t, Blob, *o.Kind())
		assert.Equal(t, int64(10), o.Size())
		assert.Equal(t, []byte("hello"), o.Content())
	})

	t.Run("invalid kind", func(t *testing.T) {
		_, err := NewObject("bad", 10, []byte("hello"))
		assert.Error(t, err)
	})
}

func TestObjectString(t *testing.T) {
	o := &Object{
		objKind: Blob,
		content: []byte("hello"),
	}

	assert.Equal(t, "hello", o.String())
}

func TestNewTreeEntry(t *testing.T) {
	obj := &Object{objKind: Blob, content: []byte("world")}
	hash, _ := HashFromString("f1d2d2f924e986ac86fdf7b36c94bcdf32beec15")

	entry := NewTreeEntry(obj, "hello.txt", 0644, hash)

	assert.Equal(t, obj, entry.Object)
	assert.Equal(t, "hello.txt", entry.Name())
	assert.Equal(t, 0644, entry.mode)
	assert.Equal(t, hash, entry.hash)
}
