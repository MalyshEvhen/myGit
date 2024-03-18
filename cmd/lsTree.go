package cmd

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
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
		if !nameOnly {
			fmt.Print("mode must be given without --name-only, and we don`t support it.")
			os.Exit(1)
		}

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
		mode, err := r.ReadString(' ')
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("read string: %w", err)
		}

		mode = mode[:len(mode)-1]
		numMode, err := strconv.Atoi(mode)
		if err != nil {
			return fmt.Errorf("atoi mode: %w", err)
		}

		name, err := r.ReadString('\000')
		if err != nil {
			return fmt.Errorf("read string: %w", err)
		}

		name = name[:len(name)-1]

		sha := make([]byte, sha1.Size)
		_, err = r.Read(sha)
		if err != nil {
			return fmt.Errorf("read sha: %w", err)
		}

		hashStr := hex.EncodeToString(sha[:])

		nestedObjHash, err := object.HashFromString(hashStr)
		if err != nil {
			return fmt.Errorf("hash from string: %w", err)
		}
		nestedObj, err := object.LoadByHash(nestedObjHash)
		if err != nil {
			return fmt.Errorf("load object: %w", err)
		}

		switch typ := nestedObj.(type) {
		case *object.Object[object.Blob]:
			fmt.Printf("%06d blob %s %s\n", numMode, hashStr, name)
		case *object.Object[object.Tree]:
			fmt.Printf("%06d tree %s %s\n", numMode, hashStr, name)
		default:
			return fmt.Errorf("unknown object type: %T", typ)
		}
	}
	return nil
}

func init() {
	lsTreeCmd.Flags().BoolVarP(&nameOnly, "name-only", "n", false, "usage: mygit ls-tree --name-only <hash>")
}
