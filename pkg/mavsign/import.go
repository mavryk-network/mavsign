package mavsign

import (
	"context"
	"fmt"

	"github.com/mavryk-network/mavbingo/v2/b58"
	"github.com/mavryk-network/mavbingo/v2/crypt"
	"github.com/mavryk-network/mavsign/pkg/utils"
	"github.com/mavryk-network/mavsign/pkg/vault"
	log "github.com/sirupsen/logrus"
)

// Import a keyPair inside the vault
func (s *MavSign) Import(ctx context.Context, importerName string, secretKey string, passCB func() ([]byte, error), opt utils.Options) (*PublicKey, error) {
	v, ok := s.vaults[importerName]
	if !ok {
		return nil, fmt.Errorf("import: vault %s is not found", importerName)
	}

	importer, ok := v.(vault.Importer)
	if !ok {
		return nil, fmt.Errorf("import: vault %s doesn't support import operation", importerName)
	}

	maybeEncrypted, err := b58.ParseEncryptedPrivateKey([]byte(secretKey))
	if err != nil {
		return nil, err
	}
	decrypted, err := maybeEncrypted.Decrypt(passCB)
	if err != nil {
		return nil, err
	}
	priv, err := crypt.NewPrivateKey(decrypted)
	if err != nil {
		return nil, err
	}
	pub := priv.Public()
	hash := pub.Hash()
	l := s.logger().WithFields(log.Fields{
		logPKH:   hash,
		logVault: importer.Name(),
	})
	l.Info("Requesting import operation")

	ref, err := importer.Import(ctx, priv, opt)
	if err != nil {
		return nil, err
	}

	s.cache.push(&keyVaultPair{pkh: hash, key: ref})

	l.WithField(logPKH, hash).Info("Successfully imported")
	pol := s.fetchPolicyOrDefault(hash)
	return &PublicKey{
		KeyReference: ref,
		Hash:         hash,
		Policy:       s.fetchPolicyOrDefault(hash),
		Active:       pol != nil,
	}, nil
}
