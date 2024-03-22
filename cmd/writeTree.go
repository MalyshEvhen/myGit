package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mygit/object"
	"github.com/spf13/cobra"
)

const (
	ModeTree = 1 << 14
	ModeBlob = 1 << 15
)

var writeTreeCmd = &cobra.Command{
	Use:   "write-tree",
	Short: "git-write-tree - Create a tree object from the current index",
	Long: `Creates a tree object using the current index.
	The name of the new tree object is printed to standard output.`,
	Run: WriteTreeCmd,
}

func WriteTreeCmd(cmd *cobra.Command, args []string) {
	hash, err := writeTree(".")
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	fmt.Println(hash)
}

func writeTree(dir string) (object.Hash, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return object.Hash{}, fmt.Errorf("read dir %v: %w", dir, err)
	}

	var table []byte

	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}

		filePath := filepath.Join(dir, f.Name())

		inf, err := f.Info()
		if err != nil {
			return object.Hash{}, fmt.Errorf("file info %v: %w", filePath, err)
		}

		mode := inf.Mode()

		if !mode.IsRegular() && !mode.IsDir() {
			continue
		}

		var hash object.Hash
		var gitMode int

		if mode.IsDir() {
			gitMode |= ModeTree

			hash, err = writeTree(filePath)
			if err != nil {
				return object.Hash{}, err
			}
		} else {
			gitMode |= ModeBlob
			gitMode |= int(mode) & 0o777

			r, size, err := object.ReadFromFile(filePath, "blob")
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			hash, err = object.Store(r, "blob", size, false)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}

		table = fmt.Appendf(table, "%o %v\000%s", gitMode, f.Name(), hash[:])
	}

	hash, err := object.Store(bytes.NewReader(table), "tree", int64(len(table)), false)
	if err != nil {
		return object.Hash{}, fmt.Errorf("store tree %v: %w", dir, err)
	}

	return hash, nil
}

func init() {
}
