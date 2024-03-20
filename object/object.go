package object

import (
	"bufio"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

type Object interface {
	String() string
	Content() []byte
	Size() int64
}

type GitObject[T GitObjectType] struct {
	content []byte
	size    int64
}

type GitObjectType interface {
	Commit | Tree | Blob | Tag
}

type Commit string
type Tree string
type Blob string
type Tag string

func (b *GitObject[T]) String() string {
	return string(b.content)
}

func (b *GitObject[T]) Size() int64 {
	return b.size
}

func (b *GitObject[T]) Content() []byte {
	return b.content
}

func LoadByHash(h Hash) (Object, error) {
	name := h.String()

	path := filepath.Join(".git", "objects", name[:2], name[2:])

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	defer func() {
		e := file.Close()
		if err == nil && e != nil {
			err = fmt.Errorf("close file: %w", e)
		}
	}()

	return LoadFile(file)
}

func LoadFile(r io.Reader) (Object, error) {
	zr, err := zlib.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("new zlib reader %w", err)
	}

	defer func() {
		e := zr.Close()
		if err == nil && e != nil {
			err = fmt.Errorf("close zlib reader: %w", e)
		}
	}()

	typ, content, err := parseObject(zr)
	if err != nil {
		return nil, fmt.Errorf("parse object %w", err)
	}

	switch typ {
	case "blob":
		return &GitObject[Blob]{content, int64(len(content))}, nil
	case "tree":
		return &GitObject[Tree]{content, int64(len(content))}, nil
	default:
		return nil, fmt.Errorf("unknown object type %s", typ)
	}
}

func NewBlob(content []byte, i int64) GitObject[Blob] {
	panic("unimplemented")
}

func parseObject(r io.Reader) (string, []byte, error) {
	br := bufio.NewReader(r)

	typ, err := br.ReadString(' ')
	if err != nil {
		return "", nil, err
	}

	typ = typ[:len(typ)-1]

	sizeStr, err := br.ReadString('\000')
	if err != nil {
		return typ, nil, err
	}

	sizeStr = sizeStr[:len(sizeStr)-1]

	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return "", nil, fmt.Errorf("parse size: %w", err)
	}

	content := make([]byte, size)

	if _, err := io.ReadFull(br, content); err != nil {
		return "", nil, fmt.Errorf("read content: %w", err)
	}

	return typ, content, nil
}
