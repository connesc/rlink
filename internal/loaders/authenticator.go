package loaders

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/sha3"

	"github.com/connesc/rlink/pkg/path"
	"github.com/connesc/rlink/pkg/path/aessiv"
	"github.com/connesc/rlink/pkg/path/mac"
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
		return mac.NewAuthenticator(mac.HMAC(md5.New), key), nil
	case "auth-hmac-sha1":
		return mac.NewAuthenticator(mac.HMAC(sha1.New), key), nil
	case "auth-hmac-sha224":
		return mac.NewAuthenticator(mac.HMAC(sha256.New224), key), nil
	case "auth-hmac-sha256":
		return mac.NewAuthenticator(mac.HMAC(sha256.New), key), nil
	case "auth-hmac-sha384":
		return mac.NewAuthenticator(mac.HMAC(sha512.New384), key), nil
	case "auth-hmac-sha512":
		return mac.NewAuthenticator(mac.HMAC(sha512.New), key), nil
	case "auth-hmac-sha512-224":
		return mac.NewAuthenticator(mac.HMAC(sha512.New512_224), key), nil
	case "auth-hmac-sha512-256":
		return mac.NewAuthenticator(mac.HMAC(sha512.New512_256), key), nil
	case "auth-sha3-224":
		return mac.NewAuthenticator(mac.PrefixedHash(sha3.New224), key), nil
	case "auth-sha3-256":
		return mac.NewAuthenticator(mac.PrefixedHash(sha3.New256), key), nil
	case "auth-sha3-384":
		return mac.NewAuthenticator(mac.PrefixedHash(sha3.New384), key), nil
	case "auth-sha3-512":
		return mac.NewAuthenticator(mac.PrefixedHash(sha3.New512), key), nil
	case "auth-shake128":
		return mac.NewAuthenticator(mac.PrefixedShakeHash(sha3.NewShake128, 32), key), nil
	case "auth-shake256":
		return mac.NewAuthenticator(mac.PrefixedShakeHash(sha3.NewShake256, 64), key), nil
	case "auth-blake2s":
		return mac.NewAuthenticator(mac.KeyedHash(blake2s.New256), key), nil
	case "auth-blake2s-128":
		return mac.NewAuthenticator(mac.KeyedHash(blake2s.New128), key), nil
	case "auth-blake2b":
		return mac.NewAuthenticator(mac.KeyedHash(blake2b.New512), key), nil
	case "auth-blake2b-256":
		return mac.NewAuthenticator(mac.KeyedHash(blake2b.New256), key), nil
	case "auth-blake2b-384":
		return mac.NewAuthenticator(mac.KeyedHash(blake2b.New384), key), nil
	case "authenc-aes-siv":
		return aessiv.NewAuthenticator(key)
	default:
		return nil, fmt.Errorf("Authenticator: unknown mode: %v", l.Mode)
	}
}
