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

	"github.com/connesc/rlink/pkg/rewriter"
)

type PathRewriter struct {
	Mode   string
	Secret Secret
}

func (l *PathRewriter) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&l.Mode, "mode", "auth-hmac-sha1", "link authentication mode")
	l.Secret.Init(cmd)
}

func (l *PathRewriter) Load() (rewriter.PathRewriter, error) {
	secret, err := l.Secret.Load()
	if err != nil {
		return nil, err
	}

	switch l.Mode {
	case "auth-hmac-md5":
		return rewriter.NewHMAC(md5.New, secret), nil
	case "auth-hmac-sha1":
		return rewriter.NewHMAC(sha1.New, secret), nil
	case "auth-hmac-sha224":
		return rewriter.NewHMAC(sha256.New224, secret), nil
	case "auth-hmac-sha256":
		return rewriter.NewHMAC(sha256.New, secret), nil
	case "auth-hmac-sha384":
		return rewriter.NewHMAC(sha512.New384, secret), nil
	case "auth-hmac-sha512":
		return rewriter.NewHMAC(sha512.New, secret), nil
	case "auth-hmac-sha512-224":
		return rewriter.NewHMAC(sha512.New512_224, secret), nil
	case "auth-hmac-sha512-256":
		return rewriter.NewHMAC(sha512.New512_256, secret), nil
	case "authenc-gcm-aes128":
		return newAESGCMRewriter(16, secret)
	case "authenc-gcm-aes192":
		return newAESGCMRewriter(24, secret)
	case "authenc-gcm-aes256":
		return newAESGCMRewriter(32, secret)
	default:
		return nil, fmt.Errorf("PathRewriter: unknown mode: %v", l.Mode)
	}
}

func newAESGCMRewriter(keySize int, secret []byte) (rewriter.PathRewriter, error) {
	if len(secret) != keySize {
		return nil, fmt.Errorf("PathRewriter: invalid AES key: expected %v bytes, got %v", keySize, len(secret))
	}
	aesCipher, err := aes.NewCipher(secret)
	if err != nil {
		return nil, fmt.Errorf("PathRewriter: cannot initialize AES cipher: %v", err)
	}
	aesGCM, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return nil, fmt.Errorf("PathRewriter: cannot initialize GCM: %v", err)
	}
	return rewriter.NewAEAD(aesGCM)
}
