package cmd

import (
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/connesc/rlink/internal/loaders"
	"github.com/connesc/rlink/pkg/proxy"
)

var proxyFlags struct {
	addr         string
	pathRewriter loaders.PathRewriter
}

var proxyCmd = &cobra.Command{
	Use:   "proxy BACKEND_URL",
	Short: "Expose a backend server through unguessable links",
	Args:  cobra.ExactArgs(1),
	Run:   runProxy,
}

func init() {
	proxyCmd.Flags().StringVar(&proxyFlags.addr, "addr", "127.0.0.1:8080", "listen address")
	proxyCmd.MarkFlagRequired("addr")
	proxyFlags.pathRewriter.Init(proxyCmd)
	rootCmd.AddCommand(proxyCmd)
}

func runProxy(cmd *cobra.Command, args []string) {
	pathRewriter, err := proxyFlags.pathRewriter.Load()
	if err != nil {
		log.Fatalln(err)
	}

	server, err := proxy.New(args[0], pathRewriter)
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(http.ListenAndServe(proxyFlags.addr, server))
}
