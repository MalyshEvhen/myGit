package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var writeTreeCmd = &cobra.Command{
	Use:   "write-tree",
	Short: "git-write-tree - Create a tree object from the current index",
	Long: `Creates a tree object using the current index.
	The name of the new tree object is printed to standard output.`,
	Run: writeTree,
}

func writeTree(cmd *cobra.Command, args []string) {
	fmt.Println("writeTree called")
}

func init() {
}
