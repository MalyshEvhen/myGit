package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mygit/object"
	"github.com/spf13/cobra"
)

var (
	message    string
	parentHash string
)

var commitTreeCmd = &cobra.Command{
	Use:   "commit-tree",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: CommitTreeCmd,
}

func CommitTreeCmd(cmd *cobra.Command, args []string) {
	if len(args) < 1 || message == "" {
		printErrAndExit()
	}
	treeHash := args[0]

	var buf []byte

	fmt.Appendf(buf, "tree %s\n", treeHash)

	if parentHash != "" {
		fmt.Appendf(buf, "parent %s\n", treeHash)
	}

	name := "Evhen"
	email := "malysh.evgeniy@gmail.com"
	now := time.Now()

	buf = fmt.Appendf(buf, "author %s <%s> %d %s\n", name, email, now.Unix(), now.Format("-0700"))
	buf = fmt.Appendf(buf, "committer %s <%s> %d %s\n", name, email, now.Unix(), now.Format("-0700"))

	buf = fmt.Appendf(buf, "\n%s", message)
	if !strings.HasPrefix(message, "\n") {
		buf = append(buf, '\n')
	}

	hash, err := object.Store(bytes.NewBuffer(buf), object.Commit, int64(len(buf)), false)
	if err != nil {
		fmt.Printf("Error %s:", err.Error())
	}

	fmt.Println(hash)

	fmt.Println("commitTree called")
}

func printErrAndExit() {
	fmt.Fprintf(os.Stderr, "usage: mygit commit-tree <tree_hash> -p <parent_hash> -m <message>\n")
	os.Exit(1)
}

func init() {
	commitTreeCmd.Flags().StringVarP(&message, "message", "m", "", "Commit message")
	commitTreeCmd.Flags().StringVarP(&parentHash, "parrent", "p", "", "Parrent commit hash")
}
