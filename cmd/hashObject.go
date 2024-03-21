package cmd

import (
	"fmt"

	"github.com/mygit/object"
	"github.com/spf13/cobra"
)

var (
	filePath string
	write    bool
)

var hashObjectCmd = &cobra.Command{
	Use:   "hash-object",
	Short: "Compute object ID and optionally create an object from a file",
	Long: `Computes the object ID value for an object with specified type with the contents of the named file (which can be outside of the work tree), and
       optionally writes the resulting object into the object database. Reports its object ID to its standard output. When <type> is not specified, it
       defaults to "blob".`,
	Run: hashObjCmd,
}

func hashObjCmd(cmd *cobra.Command, args []string) {
	if len(args) <= 0 {
		cmd.Help()
	} else {
		filePath = args[0]
		r, size, err := object.ReadFromFile(filePath, "blob")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		hash, err := object.Store(r, "blob", size, !write)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		fmt.Println(hash)
	}
}

func init() {
	hashObjectCmd.Flags().BoolVarP(&write, "write", "w", false, "write the object into the object database")
}
