package object

import (
	"bufio"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"strconv"
)

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
