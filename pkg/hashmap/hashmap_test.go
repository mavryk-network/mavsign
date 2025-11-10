package hashmap_test

import (
	"encoding/json"
	"testing"

	mv "github.com/mavryk-network/mavbingo/v2"
	"github.com/mavryk-network/mavbingo/v2/crypt"
	"github.com/mavryk-network/mavsign/pkg/hashmap"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestHashMap(t *testing.T) {
	m := hashmap.New[mv.EncodedPublicKeyHash]([]hashmap.KV[crypt.PublicKeyHash, string]{
		{
			&mv.Ed25519PublicKeyHash{0}, "a",
		},
		{
			&mv.Ed25519PublicKeyHash{1}, "b",
		},
		{
			&mv.Ed25519PublicKeyHash{2}, "c",
		},
		{
			&mv.Ed25519PublicKeyHash{3}, "d",
		},
	})

	t.Run("JSON", func(t *testing.T) {
		buf, err := json.Marshal(m)
		require.NoError(t, err)

		var res hashmap.HashMap[mv.EncodedPublicKeyHash, crypt.PublicKeyHash, string]
		err = json.Unmarshal(buf, &res)
		require.NoError(t, err)
		require.Equal(t, m, res)
	})

	t.Run("YAML", func(t *testing.T) {
		buf, err := yaml.Marshal(m)
		require.NoError(t, err)
		var res hashmap.HashMap[mv.EncodedPublicKeyHash, crypt.PublicKeyHash, string]
		err = yaml.Unmarshal(buf, &res)
		require.NoError(t, err)
		require.Equal(t, m, res)
	})
}
