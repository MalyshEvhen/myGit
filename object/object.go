package object

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Object struct {
	objKind Kind
	size    int64
	content []byte
}

type Kind string

const (
	Commit Kind = "commit"
	Tree   Kind = "tree"
	Blob   Kind = "blob"
)

func (k Kind) String() string {
	switch k {
	case Commit:
		return "commit"
	case Tree:
		return "tree"
	case Blob:
		return "blob"
	}
	return ""
}

func NewObject(kind string, size int64, content []byte) (*Object, error) {
	objKind := Kind(kind)
	if objKind != Commit && objKind != Tree && objKind != Blob {
		return nil, fmt.Errorf("unsupported git object type: %s", kind)
	}

	return &Object{objKind, size, content}, nil
}

func (o *Object) Kind() *Kind {
	return &o.objKind
}

func (o *Object) Size() int64 {
	return o.size
}

func (o *Object) Content() []byte {
	return o.content
}

func (o *Object) String() string {
	return string(o.content)
}

func ReadObject(r io.Reader) (*Object, error) {
	br := bufio.NewReader(r)

	kind, err := br.ReadString(' ')
	if err != nil {
		return nil, err
	}

	kind = kind[:len(kind)-1]

	sizeStr, err := br.ReadString('\000')
	if err != nil {
		return nil, err
	}

	sizeStr = sizeStr[:len(sizeStr)-1]

	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse size: %w", err)
	}

	content := make([]byte, size)

	if _, err := io.ReadFull(br, content); err != nil {
		return nil, fmt.Errorf("read content: %w", err)
	}

	return NewObject(kind, int64(len(content)), content)
}
