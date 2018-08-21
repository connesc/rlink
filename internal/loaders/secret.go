package loaders

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
)

type Secret struct {
	Secret         string
	SecretEncoding string
}

func (l *Secret) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&l.Secret, "secret", "", "secret pass")
	cmd.MarkFlagRequired("secret")
	cmd.Flags().StringVar(&l.SecretEncoding, "secret-encoding", "utf8", "encoding of the secret pass (\"utf8\", \"base64\" or \"hex\")")
}

func (l *Secret) Load() ([]byte, error) {
	switch l.SecretEncoding {
	case "utf8":
		return []byte(l.Secret), nil
	case "base64":
		secret, err := base64.StdEncoding.DecodeString(l.Secret)
		if err != nil {
			return nil, fmt.Errorf("Secret: failed to decode base64: %v", err)
		}
		return []byte(secret), nil
	case "hex":
		secret, err := hex.DecodeString(l.Secret)
		if err != nil {
			return nil, fmt.Errorf("Secret: failed to decode hex: %v", err)
		}
		return []byte(secret), nil
	default:
		return nil, fmt.Errorf("Secret: unknown encoding: %v", l.SecretEncoding)
	}
}
