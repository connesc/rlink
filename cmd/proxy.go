package cmd

import (
	"net/http"

	"github.com/connesc/rlink/pkg/proxy"
	"github.com/spf13/cobra"
)

var proxyFlags struct {
	addr           string
	mode           string
	secret         string
	secretEncoding string
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
	proxyCmd.Flags().StringVar(&proxyFlags.mode, "mode", "sign", "link mode (\"sign\" or \"encrypt\")")
	proxyCmd.Flags().StringVar(&proxyFlags.secret, "secret", "", "secret pass")
	proxyCmd.MarkFlagRequired("secret")
	proxyCmd.Flags().StringVar(&proxyFlags.secretEncoding, "secret-encoding", "utf8", "encoding of the secret pass (\"utf8\" or \"base64\")")
	rootCmd.AddCommand(proxyCmd)
}

func runProxy(cmd *cobra.Command, args []string) {
	server := proxy.New(args[0], []byte(proxyFlags.secret))
	panic(http.ListenAndServe(proxyFlags.addr, server))
}
