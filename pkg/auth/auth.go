package auth

import (
	"context"
	stderr "errors"
	"net/http"

	"github.com/mavryk-network/mavbingo/v2"
	"github.com/mavryk-network/mavbingo/v2/crypt"
	"github.com/mavryk-network/mavsign/pkg/errors"
)

// ErrPublicKeyNotFound is returned by AuthorizedKeysStorage.GetPublicKey if authorized key is not found
var ErrPublicKeyNotFound = errors.Wrap(stderr.New("public key not found"), http.StatusUnauthorized)

// AuthorizedKeysStorage represents an authorized public keys storage
type AuthorizedKeysStorage interface {
	GetPublicKey(ctx context.Context, keyHash mavbingo.PublicKeyHash) (crypt.PublicKey, error)
	ListPublicKeys(ctx context.Context) ([]crypt.PublicKeyHash, error)
}

// Must panics in case of error
func Must(s AuthorizedKeysStorage, err error) AuthorizedKeysStorage {
	if err != nil {
		panic(err)
	}
	return s
}
