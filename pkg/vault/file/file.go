package file

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	config "github.com/mavryk-network/mavsign/pkg/config"
	"github.com/mavryk-network/mavsign/pkg/errors"
	"github.com/mavryk-network/mavsign/pkg/vault"
	"github.com/mavryk-network/mavsign/pkg/vault/memory"
	"gopkg.in/yaml.v3"
)

type cfg struct {
	File string `yaml:"file" validate:"required"`
}

func trimSecretKey(k string) string {
	if i := strings.IndexByte(string(k), ':'); i >= 0 {
		return string(k)[i+1:]
	}
	return string(k)
}

type mavrykSecretJSONEntry struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func init() {
	vault.RegisterVault("file", func(ctx context.Context, node *yaml.Node) (vault.Vault, error) {
		var conf cfg
		if node == nil || node.Kind == 0 {
			return nil, errors.New("(File): config is missing")
		}
		if err := node.Decode(&conf); err != nil {
			return nil, err
		}

		if err := config.Validator().Struct(&conf); err != nil {
			return nil, err
		}

		content, err := ioutil.ReadFile(conf.File)
		if err != nil {
			return nil, fmt.Errorf("(File): %w", err)
		}

		var entries []*mavrykSecretJSONEntry
		if err := json.Unmarshal(content, &entries); err != nil {
			return nil, fmt.Errorf("(File): %w", err)
		}

		data := make([]*memory.UnparsedKey, len(entries))
		for i, e := range entries {
			data[i] = &memory.UnparsedKey{
				Data: trimSecretKey(e.Value),
				ID:   e.Name,
			}
		}
		return memory.NewUnparsed(data, "File"), nil
	})
}
