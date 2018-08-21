package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/connesc/rlink/internal/loaders"
)

var indexFlags struct {
	addr         string
	pathRewriter loaders.PathRewriter
	parent       bool
}

var indexCmd = &cobra.Command{
	Use:   "index ROOT_PATH SERVER_URL",
	Short: "Provide an index for files exposed by a rlink server",
	Args:  cobra.ExactArgs(2),
	Run:   runIndex,
}

func init() {
	indexCmd.Flags().StringVar(&indexFlags.addr, "addr", "127.0.0.1:8080", "listen address")
	indexCmd.MarkFlagRequired("addr")
	indexFlags.pathRewriter.Init(indexCmd)
	indexCmd.Flags().BoolVar(&indexFlags.parent, "parent", false, "whether to link to the parent directory")
	rootCmd.AddCommand(indexCmd)
}

func runIndex(cmd *cobra.Command, args []string) {
	fmt.Println("index", indexFlags.addr, args[0], args[1])
}
