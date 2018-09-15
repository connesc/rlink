package loaders

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/connesc/rlink/pkg/path"
	"github.com/connesc/rlink/pkg/path/aead"
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
	case "authenc-gcm-aes128":
		return newAESGCMRewriter(16, key)
	case "authenc-gcm-aes192":
		return newAESGCMRewriter(24, key)
	case "authenc-gcm-aes256":
		return newAESGCMRewriter(32, key)
	default:
		return nil, fmt.Errorf("Authenticator: unknown mode: %v", l.Mode)
	}
}

func newAESGCMRewriter(keySize int, key []byte) (path.Authenticator, error) {
	if len(key) != keySize {
		return nil, fmt.Errorf("Authenticator: invalid AES key: expected %v bytes, got %v", keySize, len(key))
	}
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("Authenticator: cannot initialize AES cipher: %v", err)
	}
	aesGCM, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return nil, fmt.Errorf("Authenticator: cannot initialize GCM: %v", err)
	}
	return aead.NewAuthenticator(aesGCM)
}
