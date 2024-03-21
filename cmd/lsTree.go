package cmd

import (
	"bufio"
	"bytes"
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

	obj, err := object.Read(hash)
	if err != nil {
		return fmt.Errorf("load object: %w", err)
	}
	if *obj.Kind() != object.Tree {
		return fmt.Errorf("object `%s` is not a tree", string(*obj.Kind()))
	}

	r := bufio.NewReader(bytes.NewReader(obj.Content()))

	for {
		entry, err := object.LoadTreeEntry(r)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		if !nameOnly {
			fmt.Printf("%s", entry)
		} else {
			fmt.Println(entry.Name())
		}
	}
	return nil
}

func init() {
	lsTreeCmd.Flags().BoolVarP(&nameOnly, "name-only", "n", false, "usage: mygit ls-tree --name-only <hash>")
}
