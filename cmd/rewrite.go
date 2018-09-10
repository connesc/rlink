package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/connesc/rlink/internal/loaders"
)

var rewriteFlags struct {
	pathRewriter loaders.PathRewriter
	reverse      bool
}

var rewriteCmd = &cobra.Command{
	Use:   "rewrite PATH",
	Short: "Rewrite the given absolute path",
	Args:  cobra.ExactArgs(1),
	Run:   runRewrite,
}

func init() {
	rewriteFlags.pathRewriter.Init(rewriteCmd)
	rewriteCmd.Flags().BoolVarP(&rewriteFlags.reverse, "reverse", "r", false, "retrieve the original path by applying the reverse transformation")
	rootCmd.AddCommand(rewriteCmd)
}

func runRewrite(cmd *cobra.Command, args []string) {
	pathRewriter, err := rewriteFlags.pathRewriter.Load()
	if err != nil {
		log.Fatalln(err)
	}

	var transform func(string) (string, error)
	if rewriteFlags.reverse {
		transform = pathRewriter.ToOriginal
	} else {
		transform = pathRewriter.FromOriginal
	}

	outputPath, err := transform(args[0])
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(outputPath)
}
