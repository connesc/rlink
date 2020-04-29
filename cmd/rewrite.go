package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/connesc/rlink/internal/loaders"
	"github.com/connesc/rlink/pkg/path"
)

var rewriteFlags struct {
	authenticator loaders.Authenticator
	reverse       bool
}

var rewriteCmd = &cobra.Command{
	Use:   "rewrite PATH",
	Short: "Rewrite the given path",
	Args:  cobra.ExactArgs(1),
	Run:   runRewrite,
}

func init() {
	rewriteFlags.authenticator.Init(rewriteCmd)
	rewriteCmd.Flags().BoolVarP(&rewriteFlags.reverse, "reverse", "r", false, "retrieve the original path by applying the reverse transformation")
	rootCmd.AddCommand(rewriteCmd)
}

func runRewrite(cmd *cobra.Command, args []string) {
	authenticator, err := rewriteFlags.authenticator.Load()
	if err != nil {
		log.Fatalln(err)
	}

	transform := authenticator.FromOriginal
	if rewriteFlags.reverse {
		transform = authenticator.ToOriginal
	}

	outputPath, err := transform(path.Normalize(args[0]))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(outputPath)
}
