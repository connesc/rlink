package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/connesc/rlink/internal/loaders"
	"github.com/connesc/rlink/pkg/server"
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
	indexCmd.Flags().BoolVar(&indexFlags.parent, "parent", true, "whether to link to the parent directory")
	rootCmd.AddCommand(indexCmd)
}

func runIndex(cmd *cobra.Command, args []string) {
	// TODO: handle args[1] (server URL)

	pathRewriter, err := indexFlags.pathRewriter.Load()
	if err != nil {
		log.Fatalln(err)
	}

	handler, err := server.New(args[0], pathRewriter, &server.Options{
		Files:       false,
		Index:       true,
		IndexParent: indexFlags.parent,
	})
	if err != nil {
		log.Fatalln(err)
	}

	root, err := pathRewriter.FromOriginal("/")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Root: http://%s/%s\n", indexFlags.addr, root)

	log.Fatalln(http.ListenAndServe(indexFlags.addr, handler))
}
