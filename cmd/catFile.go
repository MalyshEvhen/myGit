package cmd

import (
	"fmt"
	"os"

	"github.com/mygit/object"
	"github.com/spf13/cobra"
)

var (
	objectHash  string
	prettyPrint bool
)

var catFileCmd = &cobra.Command{
	Use:   "cat-file",
	Short: "Provide contents or details of repository objects",
	Long: `Output the contents or other properties such as size, type or delta information of one or more objects.
       This command can operate in two modes, depending on whether an option from the --batch family is specified.
       In non-batch mode, the command provides information on an object named on the command line.
       In batch mode, arguments are read from standard input.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) <= 0 {
			cmd.Help()
		} else {
			objectHash = args[0]
			if !prettyPrint {
				fmt.Print("mode must be given without -p, and we don`t support it.")
				os.Exit(1)
			}
			if err := catFile(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
	},
}

func catFile() error {
	hash, err := object.HashFromString(objectHash)
	if err != nil {
		return fmt.Errorf("hash object: %w, %v", err, objectHash)
	}

	obj, err := object.LoadByHash(hash)
	if err != nil {
		return fmt.Errorf("load object: %w", err)
	}
	switch obj.(type) {
	case *object.Object[object.Blob]:
		fmt.Printf("%s", obj.Content())
	case *object.Object[object.Tree]:
		fmt.Printf("tree %s\n", objectHash)
	default:
		fmt.Printf("unknown type of object: %v", obj)
		os.Exit(1)
	}
	return nil
}

func init() {
	catFileCmd.Flags().BoolVarP(&prettyPrint, "pretty", "p", false, "Pretty print the contents of the object")
}
