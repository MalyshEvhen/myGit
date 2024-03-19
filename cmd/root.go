package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var objectHash string

var rootCmd = &cobra.Command{
	Use:   "mygit",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
