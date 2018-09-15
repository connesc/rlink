package cmd

import (
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/connesc/rlink/internal/loaders"
	"github.com/connesc/rlink/pkg/proxy"
)

var proxyFlags struct {
	addr          string
	authenticator loaders.Authenticator
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
	proxyFlags.authenticator.Init(proxyCmd)
	rootCmd.AddCommand(proxyCmd)
}

func runProxy(cmd *cobra.Command, args []string) {
	authenticator, err := proxyFlags.authenticator.Load()
	if err != nil {
		log.Fatalln(err)
	}

	handler, err := proxy.New(args[0], authenticator)
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(http.ListenAndServe(proxyFlags.addr, handler))
}
