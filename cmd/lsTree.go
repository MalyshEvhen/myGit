package cmd

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"os"

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
		_, err := r.ReadString(' ')
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("read string: %w", err)
		}

		name, err := r.ReadString('\000')
		if err != nil {
			return fmt.Errorf("read string: %w", err)
		}

		name = name[:len(name)-1]

		_, err = r.Discard(sha1.Size)
		if err != nil {
			return fmt.Errorf("discard: %w", err)
		}

		fmt.Printf("%s\n", name)
	}
	return nil
}

func init() {
	lsTreeCmd.Flags().BoolVarP(&nameOnly, "name-only", "n", false, "usage: mygit ls-tree --name-only <hash>")
}
