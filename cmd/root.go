package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var objectHash string

var rootCmd = &cobra.Command{
	Use:   "mygit",
	Short: "A CLI tool to interact with a git repository",
	Long:  `A CLI tool to interact with a git repository`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(catFileCmd)
	rootCmd.AddCommand(hashObjectCmd)
	rootCmd.AddCommand(lsTreeCmd)
	rootCmd.AddCommand(writeTreeCmd)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
