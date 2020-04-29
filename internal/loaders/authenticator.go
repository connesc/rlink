package loaders

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"regexp"
	"strconv"

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

var variableModePattern = regexp.MustCompile(`^(.*)-([0-9]+)$`)

func parseMode(mode string) (string, int) {
	matches := variableModePattern.FindStringSubmatch(mode)
	if matches == nil {
		switch mode {
		case "auth-shake128":
			return mode, 32
		case "auth-shake256":
			return mode, 64
		case "auth-blake2b":
			return mode, blake2b.Size
		default:
			return mode, 0
		}
	}

	variableMode := matches[1]

	bits, err := strconv.Atoi(matches[2])
	if err != nil || bits%8 != 0 {
		return mode, 0
	}
	size := bits / 8

	switch variableMode {
	case "auth-shake128":
		if size >= 16 {
			return variableMode, size
		}
	case "auth-shake256":
		if size >= 16 {
			return variableMode, size
		}
	case "auth-blake2b":
		if size >= 16 && size <= blake2b.Size {
			return variableMode, size
		}
	}

	return mode, 0
}

func (l *Authenticator) Load() (path.Authenticator, error) {
	key, err := l.Key.Load()
	if err != nil {
		return nil, err
	}

	mode, size := parseMode(l.Mode)

	switch mode {
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
		return mac.NewAuthenticator(mac.PrefixedShakeHash(sha3.NewShake128, size), key), nil
	case "auth-shake256":
		return mac.NewAuthenticator(mac.PrefixedShakeHash(sha3.NewShake256, size), key), nil
	case "auth-blake2s":
		return mac.NewAuthenticator(mac.KeyedHash(blake2s.New256), key), nil
	case "auth-blake2s-128":
		return mac.NewAuthenticator(mac.KeyedHash(blake2s.New128), key), nil
	case "auth-blake2b":
		return mac.NewAuthenticator(mac.VariableKeyedHash(blake2b.New, size), key), nil
	case "authenc-aes-siv":
		return aessiv.NewAuthenticator(key)
	default:
		return nil, fmt.Errorf("Authenticator: unknown mode: %v", mode)
	}
}
