package loaders

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
)

type Key struct {
	Key         string
	KeyEncoding string
}

func (l *Key) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&l.Key, "key", "", "secret key")
	cmd.MarkFlagRequired("key")
	cmd.Flags().StringVar(&l.KeyEncoding, "key-encoding", "utf8", "encoding of the secret key (\"utf8\", \"base64\" or \"hex\")")
}

func (l *Key) Load() ([]byte, error) {
	switch l.KeyEncoding {
	case "utf8":
		return []byte(l.Key), nil
	case "base64":
		key, err := base64.StdEncoding.DecodeString(l.Key)
		if err != nil {
			return nil, fmt.Errorf("Key: failed to decode base64: %v", err)
		}
		return []byte(key), nil
	case "hex":
		key, err := hex.DecodeString(l.Key)
		if err != nil {
			return nil, fmt.Errorf("Key: failed to decode hex: %v", err)
		}
		return []byte(key), nil
	default:
		return nil, fmt.Errorf("Key: unknown encoding: %v", l.KeyEncoding)
	}
}
