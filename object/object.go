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

type GitObject struct {
	objType GitObjectType
	size    int64
	content []byte
}

type GitObjectType string

const (
	Commit GitObjectType = "commit"
	Tree   GitObjectType = "tree"
	Blob   GitObjectType = "blob"
)

func NewGitObject(typ string, size int64, content []byte) (*GitObject, error) {
	var objType GitObjectType
	switch typ {
	case string(Commit):
		objType = Commit
	case string(Tree):
		objType = Tree
	case string(Blob):
		objType = Blob
	default:
		return nil, fmt.Errorf("unsupported git object type: %s", typ)
	}

	return &GitObject{objType, size, content}, nil
}

func (o *GitObject) Type() *GitObjectType {
	return &o.objType
}

func (o *GitObject) Size() int64 {
	return o.size
}

func (o *GitObject) Content() []byte {
	return o.content
}

func (o *GitObject) String() string {
	return string(o.content)
}

func LoadByHash(h Hash) (*GitObject, error) {
	name := h.String()

	path := filepath.Join(".git", "objects", name[:2], name[2:])

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	defer file.Close()

	zr, err := Decompress(file)
	if err != nil {
		return nil, err
	}

	typ, content, err := parse(zr)
	if err != nil {
		return nil, fmt.Errorf("parse object %w", err)
	}

	return NewGitObject(typ, int64(len(content)), content)
}

type TreeEntry struct {
	*GitObject
	name string
	mode int
	hash Hash
}

func NewTreeEntry(obj *GitObject, name string, mode int, sha []byte) *TreeEntry {
	return &TreeEntry{
		GitObject: obj,
		name:      name,
		mode:      mode,
		hash:      Hash(sha[:]),
	}
}

func (t *TreeEntry) String() string {
	objType := t.GitObject.Type()
	objTypeValue := string(*objType)

	return fmt.Sprintf("%06d %s %s    %s\n", t.mode, objTypeValue, t.hash, t.name)
}

func Decompress(r io.Reader) (io.ReadCloser, error) {
	zr, err := zlib.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("new zlib reader %w", err)
	}
	defer zr.Close()

	return zr, nil
}

func parse(r io.Reader) (string, []byte, error) {
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
