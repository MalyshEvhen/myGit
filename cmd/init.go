package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create an empty Git repository or reinitialize an existing one.",
	Long: `For the purposes of this challenge, it only create the bare
minimum required for Git to function properly:

    - .git/
    - objects/
    - refs/
    - HEAD (contain "ref: refs/heads/master\n")`,
	Run: initialize,
}

func initialize(cmd *cobra.Command, args []string) {
	for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
		}
	}

	headFileContents := []byte("ref: refs/heads/main\n")
	if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
	}

	fmt.Println("Initialized git repository")
}
