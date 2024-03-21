package object

import (
	"bufio"
	"crypto/sha1"
	"errors"
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

type TreeEntry struct {
	*Object
	name string
	mode int
	hash Hash
}

func NewTreeEntry(o *Object, name string, mode int, hash Hash) *TreeEntry {
	return &TreeEntry{
		Object: o,
		name:   name,
		mode:   mode,
		hash:   hash,
	}
}

func (e *TreeEntry) Name() string {
	return e.name
}

func (t *TreeEntry) String() string {
	return fmt.Sprintf("%06d %s %s    %s\n", t.mode, t.Object.Kind(), t.hash, t.name)
}

func ReadTreeEntry(r *bufio.Reader) (*TreeEntry, error) {
	mode, err := readFileMode(r)
	if err != nil {
		return nil, err
	}

	name, err := readName(r)
	if err != nil {
		return nil, err
	}

	sha, err := readSha(r)
	if err != nil {
		return nil, err
	}
	hash := Hash(sha[:])

	rc, err := LoadFileByHash(hash)
	if err != nil {
		return nil, fmt.Errorf("load tree entry by hash: %s %w", hash, err)
	}
	defer rc.Close()

	obj, err := ReadObject(rc)
	if err != nil {
		return nil, err
	}

	return NewTreeEntry(obj, name, mode, hash), nil
}

func readSha(r *bufio.Reader) ([]byte, error) {
	sha := make([]byte, sha1.Size)
	_, err := r.Read(sha)
	if err != nil {
		return nil, fmt.Errorf("read sha: %w", err)
	}
	return sha, nil
}

func readName(r *bufio.Reader) (string, error) {
	name, err := r.ReadString('\000')
	if err != nil {
		return "", fmt.Errorf("read string: %w", err)
	}
	name = name[:len(name)-1]

	return name, nil
}

func readFileMode(r *bufio.Reader) (int, error) {
	mode, err := r.ReadString(' ')
	if errors.Is(err, io.EOF) {
		return 0, err
	}
	if err != nil {
		return 0, fmt.Errorf("read string: %w", err)
	}
	mode = mode[:len(mode)-1]

	modeNum, err := strconv.Atoi(mode)
	if err != nil {
		return 0, fmt.Errorf("atoi mode: %w", err)
	}
	return modeNum, nil
}
