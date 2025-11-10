package integrationtest

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGCPVault(t *testing.T) {

	project := os.Getenv("VAULT_GCP_PROJECTID")
	location := os.Getenv("VAULT_GCP_LOCATION")
	keyring := os.Getenv("VAULT_GCP_KEYRING")
	mv3 := os.Getenv("VAULT_GCP_MV3")
	mv3pk := os.Getenv("VAULT_GCP_MV3_PK")
	mv3alias := "gcptz3"

	//config
	var c Config
	c.Read()
	var v VaultConfig
	v.Driver = "cloudkms"
	v.Conf = map[string]interface{}{"project": &project, "location": &location, "key_ring": &keyring}
	c.Vaults["gcp"] = &v
	var p MavrykPolicy
	p.LogPayloads = true
	p.Allow = map[string][]string{"generic": {"reveal", "transaction"}}
	c.Mavryk[mv3] = &p
	backup_then_update_config(c)
	defer restore_config()
	restart_mavsign()

	//setup
	out, err := MavkitClient("import", "secret", "key", mv3alias, "http://mavsign:6732/"+mv3)
	assert.NoError(t, err)
	assert.Contains(t, string(out), "Mavryk address added: "+mv3)
	defer MavkitClient("forget", "address", mv3alias, "--force")

	out, err = MavkitClient("transfer", "100", "from", "alice", "to", mv3alias, "--burn-cap", "0.06425")
	assert.NoError(t, err)
	require.Contains(t, string(out), "Operation successfully injected in the node")

	//test
	out, err = MavkitClient("transfer", "1", "from", mv3alias, "to", "alice", "--burn-cap", "0.06425")
	assert.NoError(t, err)
	require.Contains(t, string(out), "Operation successfully injected in the node")

	require.Equal(t, mv3pk, GetPublicKey(mv3))
}
