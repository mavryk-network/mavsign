package file

import (
	"context"

	"github.com/mavryk-network/mavsign/pkg/errors"
	"github.com/mavryk-network/mavsign/pkg/vault"
	"github.com/mavryk-network/mavsign/pkg/vault/memory"
	"gopkg.in/yaml.v3"
)

func init() {
	vault.RegisterVault("mem", func(ctx context.Context, node *yaml.Node) (vault.Vault, error) {
		var conf []string
		if node == nil || node.Kind == 0 {
			return nil, errors.New("(Mem): config is missing")
		}
		if err := node.Decode(&conf); err != nil {
			return nil, err
		}

		data := make([]*memory.UnparsedKey, len(conf))
		for i, v := range conf {
			data[i] = &memory.UnparsedKey{Data: v}
		}
		return &memory.Importer{Vault: memory.NewUnparsed(data, "Mem")}, nil
	})
}
