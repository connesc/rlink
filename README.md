# Rlink

```
Rlink allows to share files using secret links.

Secret links are derived from a secret key using cryptographic operations. Since those operations
are deterministic, there is no need for persistent storage.

The authentication mode is generally specified with the --mode flag and allows to define which
cryptographic operation is used. Modes are grouped in two categories:
- auth-xxx: links are authenticated and paths are not encrypted.
- authenc-xxx: links are authenticated and directory names are encrypted.

In both cases, a secret link cannot be guessed. Both categories differ however on which information
is revealed when sharing a secret link to a third party.

The following modes are available:
- auth-hmac-md5:
    Authentication using a HMAC + MD5.

- auth-hmac-sha1:
    Authentication using a HMAC + SHA-1.

- auth-hmac-sha224,
  auth-hmac-sha256,
  auth-hmac-sha384,
  auth-hmac-sha512,
  auth-hmac-sha512-224,
  auth-hmac-sha512-256:
    Authentication using a HMAC + SHA-2 hash functions.

- auth-sha3-224,
  auth-sha3-256,
  auth-sha3-384,
  auth-sha3-512:
    Authentication using SHA-3 hash functions.

- auth-shake128,
  auth-shake128-xxx:
    Authentication using SHAKE-128. A custom hash length can be specified as a multiple of 8 bits.
    It must be >= 128 and defaults to 256. Consider using a multiple of 24 for optimal use of
    base64 encoding in secret links.

- auth-shake256,
  auth-shake256-xxx:
    Authentication using SHAKE-256. A custom hash length can be specified as a multiple of 8 bits.
    It must be >= 128 and defaults to 512. Consider using a multiple of 24 for optimal use of
    base64 encoding in secret links.

- auth-blake2s,
  auth-blake2s-128:
    Authentication using BLAKE2s (either 256 or 128 hash length).

- auth-blake2b,
  auth-blake2b-xxx:
    Authentication using BLAKE2b. A custom hash length can be specified as a multiple of 8 bits.
    It must be between 128 and 512 and defaults to 512. Consider using a multiple of 24 for optimal
    use of base64 encoding in secret links.

- authenc-aes-siv:
    Authenticated encryption using AES-GCM-SIV (aka. AES-SIV). Key must be 64 bytes long.

Usage:
  rlink [command]

Available Commands:
  help        Help about any command
  index       Provide an index for files exposed by a rlink server
  proxy       Expose a backend server using secret links
  rewrite     Rewrite the given path
  server      Serve files using secret links

Flags:
  -h, --help   help for rlink

Use "rlink [command] --help" for more information about a command.
```
