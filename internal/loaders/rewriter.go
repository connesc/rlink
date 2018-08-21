package loaders

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/connesc/rlink/pkg/rewriter"
)

type PathRewriter struct {
	Mode   string
	Secret Secret
}

func (l *PathRewriter) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&l.Mode, "mode", "sign", "link mode (\"sign\" or \"encrypt\")")
	l.Secret.Init(cmd)
}

func (l *PathRewriter) Load() (rewriter.PathRewriter, error) {
	secret, err := l.Secret.Load()
	if err != nil {
		return nil, err
	}

	switch l.Mode {
	case "sign":
		return rewriter.NewPathSigner(sha1.New, secret, base64.RawURLEncoding), nil
	default:
		return nil, fmt.Errorf("PathRewriter: unknown mode: %v", l.Mode)
	}
}
