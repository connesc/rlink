package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var serverFlags struct {
	addr           string
	mode           string
	secret         string
	secretEncoding string
	index          bool
	indexParent    bool
	prefix         bool
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
	serverCmd.Flags().StringVar(&serverFlags.mode, "mode", "sign", "link mode (\"sign\" or \"encrypt\")")
	serverCmd.Flags().StringVar(&serverFlags.secret, "secret", "", "secret pass")
	serverCmd.MarkFlagRequired("secret")
	serverCmd.Flags().StringVar(&serverFlags.secretEncoding, "secret-encoding", "utf8", "encoding of the secret pass (\"utf8\" or \"base64\")")
	serverCmd.Flags().BoolVar(&serverFlags.index, "index", false, "whether to provide indices for directories")
	serverCmd.Flags().BoolVar(&serverFlags.indexParent, "index-parent", false, "whether to link to the parent directory in indices")
	serverCmd.Flags().BoolVar(&serverFlags.prefix, "prefix", false, "whether to prefix links with /file and /dir")
	rootCmd.AddCommand(serverCmd)
}

func runServer(cmd *cobra.Command, args []string) {
	fmt.Println("server", serverFlags.addr, args[0])
}
