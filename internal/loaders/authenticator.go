package loaders

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/connesc/rlink/pkg/path"
	"github.com/connesc/rlink/pkg/path/aessiv"
	"github.com/connesc/rlink/pkg/path/hmac"
)

type Authenticator struct {
	Mode string
	Key  Key
}

func (l *Authenticator) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&l.Mode, "mode", "auth-hmac-sha1", "link authentication mode")
	l.Key.Init(cmd)
}

func (l *Authenticator) Load() (path.Authenticator, error) {
	key, err := l.Key.Load()
	if err != nil {
		return nil, err
	}

	switch l.Mode {
	case "auth-hmac-md5":
		return hmac.NewAuthenticator(md5.New, key), nil
	case "auth-hmac-sha1":
		return hmac.NewAuthenticator(sha1.New, key), nil
	case "auth-hmac-sha224":
		return hmac.NewAuthenticator(sha256.New224, key), nil
	case "auth-hmac-sha256":
		return hmac.NewAuthenticator(sha256.New, key), nil
	case "auth-hmac-sha384":
		return hmac.NewAuthenticator(sha512.New384, key), nil
	case "auth-hmac-sha512":
		return hmac.NewAuthenticator(sha512.New, key), nil
	case "auth-hmac-sha512-224":
		return hmac.NewAuthenticator(sha512.New512_224, key), nil
	case "auth-hmac-sha512-256":
		return hmac.NewAuthenticator(sha512.New512_256, key), nil
	case "authenc-aes-siv":
		return aessiv.NewAuthenticator(key)
	default:
		return nil, fmt.Errorf("Authenticator: unknown mode: %v", l.Mode)
	}
}
