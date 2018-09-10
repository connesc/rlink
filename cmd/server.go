package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/connesc/rlink/internal/loaders"
	"github.com/connesc/rlink/pkg/server"
)

var serverFlags struct {
	addr         string
	pathRewriter loaders.PathRewriter
	index        bool
	indexParent  bool
	prefix       bool
}

var serverCmd = &cobra.Command{
	Use:   "server ROOT_PATH",
	Short: "Serve files through unguessable links",
	Args:  cobra.ExactArgs(1),
	Run:   runServer,
}

func init() {
	serverCmd.Flags().StringVar(&serverFlags.addr, "addr", "127.0.0.1:8080", "listen address")
	serverCmd.MarkFlagRequired("addr")
	serverFlags.pathRewriter.Init(serverCmd)
	serverCmd.Flags().BoolVar(&serverFlags.index, "index", false, "whether to provide indices for directories")
	serverCmd.Flags().BoolVar(&serverFlags.indexParent, "index-parent", false, "whether to link to the parent directory in indices")
	serverCmd.Flags().BoolVar(&serverFlags.prefix, "prefix", false, "whether to prefix links with /file and /dir")
	rootCmd.AddCommand(serverCmd)
}

func runServer(cmd *cobra.Command, args []string) {
	pathRewriter, err := serverFlags.pathRewriter.Load()
	if err != nil {
		log.Fatalln(err)
	}

	handler, err := server.New(args[0], pathRewriter)
	if err != nil {
		log.Fatalln(err)
	}

	if serverFlags.index {
		root, err := pathRewriter.FromOriginal("/")
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Root: http://%s/%s\n", serverFlags.addr, root)
	}

	log.Fatalln(http.ListenAndServe(serverFlags.addr, handler))
}
