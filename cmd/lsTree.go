package cmd

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/mygit/object"
	"github.com/spf13/cobra"
)

var nameOnly bool

var lsTreeCmd = &cobra.Command{
	Use:   "ls-tree",
	Short: "List the contents of a tree object",
	Long: `Lists the contents of a given tree object,
	like what "/bin/ls -a" does in the current working directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) <= 0 {
			cmd.Help()
			os.Exit(1)
		}
		objectHash = args[0]

		err := lsTree()
		if err != nil {
			fmt.Printf("error while list: %v", err)
			os.Exit(1)
		}
	},
}

func lsTree() error {
	hash, err := object.HashFromString(objectHash)
	if err != nil {
		return fmt.Errorf("hash object: %w %v", err, objectHash)
	}

	obj, err := object.LoadByHash(hash)
	if err != nil {
		return fmt.Errorf("load object: %w", err)
	}

	r := bufio.NewReader(bytes.NewReader(obj.Content()))

	for {
		mode, err := readFileMode(r)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		name, err := readName(r)
		if err != nil {
			return err
		}

		if !nameOnly {
			sha, err := readObjSha(r)
			if err != nil {
				return err
			}

			nestedObj, err := readNestedObject(sha)
			if err != nil {
				return err
			}

			entry := NewTreeEntry(nestedObj, name, mode, sha)
			fmt.Printf("%s", entry)
		} else {
			if _, err := r.Discard(sha1.Size); err != nil {
				return fmt.Errorf("discard sha: %w", err)
			}
			fmt.Println(name)
		}
	}
	return nil
}

func readNestedObject(sha []byte) (object.GitObject, error) {
	hash := object.Hash(sha[:])
	obj, err := object.LoadByHash(hash)
	if err != nil {
		return nil, fmt.Errorf("load object: %w", err)
	}

	return obj, nil
}

func readObjSha(r *bufio.Reader) ([]byte, error) {
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

type treeEntry struct {
	object.GitObject
	name string
	mode int
	hash object.Hash
}

func NewTreeEntry(gio object.GitObject, name string, mode int, sha []byte) *treeEntry {
	return &treeEntry{
		GitObject: gio,
		name:      name,
		mode:      mode,
		hash:      object.Hash(sha[:]),
	}
}

func (t *treeEntry) String() string {
	var typ string
	switch t.GitObject.(type) {
	case *object.Object[object.Blob]:
		typ = "blob"
	case *object.Object[object.Tree]:
		typ = "tree"
	default:
		return ""
	}
	return fmt.Sprintf("%06d %s %s    %s\n", t.mode, typ, t.hash, t.name)
}

func init() {
	lsTreeCmd.Flags().BoolVarP(&nameOnly, "name-only", "n", false, "usage: mygit ls-tree --name-only <hash>")
}
