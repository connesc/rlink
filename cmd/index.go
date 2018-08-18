package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var indexFlags struct {
	addr           string
	mode           string
	secret         string
	secretEncoding string
	parent         bool
}

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Provide an index for files exposed by a rlink server",
	Args:  cobra.ExactArgs(2),
	Run:   runIndex,
}

func init() {
	indexCmd.Flags().StringVar(&indexFlags.addr, "addr", "127.0.0.1:8080", "listen address")
	indexCmd.MarkFlagRequired("addr")
	indexCmd.Flags().StringVar(&indexFlags.mode, "mode", "sign", "link mode (\"sign\" or \"encrypt\")")
	indexCmd.Flags().StringVar(&indexFlags.secret, "secret", "", "secret pass")
	indexCmd.MarkFlagRequired("secret")
	indexCmd.Flags().StringVar(&indexFlags.secretEncoding, "secret-encoding", "utf8", "encoding of the secret pass (\"utf8\" or \"base64\")")
	indexCmd.Flags().BoolVar(&indexFlags.parent, "parent", false, "whether to link to the parent directory")
	rootCmd.AddCommand(indexCmd)
}

func runIndex(cmd *cobra.Command, args []string) {
	fmt.Println("index", indexFlags.addr, args[0], args[1])
}
