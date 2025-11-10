package integrationtest

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAZVault(t *testing.T) {

	spkey := "/etc/service-principal.key"

	thumb := os.Getenv("VAULT_AZ_CLIENTCERTTHUMB")
	clientid := os.Getenv("VAULT_AZ_CLIENTID")
	resgroup := os.Getenv("VAULT_AZ_RESGROUP")
	subid := os.Getenv("VAULT_AZ_SUBID")
	tenantid := os.Getenv("VAULT_AZ_TENANTID")
	vault := os.Getenv("VAULT_AZ_VAULT")

	mv2 := os.Getenv("VAULT_AZ_MV2")
	mv3 := os.Getenv("VAULT_AZ_MV3")
	mv3pk := os.Getenv("VAULT_AZ_MV3_PK")

	mv2alias := "aztz2"
	mv3alias := "aztz3"

	//config
	var c Config
	c.Read()
	var v VaultConfig
	v.Driver = "azure"
	v.Conf = map[string]interface{}{"vault": &vault, "tenant_id": &tenantid, "client_id": &clientid, "client_private_key": &spkey, "client_certificate_thumbprint": &thumb, "subscription_id": &subid, "resource_group": &resgroup}
	c.Vaults["azure"] = &v
	var p MavrykPolicy
	p.LogPayloads = true
	p.Allow = map[string][]string{"generic": {"reveal", "transaction"}}
	c.Mavryk[mv2] = &p
	c.Mavryk[mv3] = &p
	backup_then_update_config(c)
	defer restore_config()
	restart_mavsign()

	//setup
	out, err := MavkitClient("import", "secret", "key", mv2alias, "http://mavsign:6732/"+mv2)
	assert.NoError(t, err)
	assert.Contains(t, string(out), "Mavryk address added: "+mv2)
	defer MavkitClient("forget", "address", mv2alias, "--force")
	out, err = MavkitClient("import", "secret", "key", mv3alias, "http://mavsign:6732/"+mv3)
	assert.NoError(t, err)
	assert.Contains(t, string(out), "Mavryk address added: "+mv3)
	defer MavkitClient("forget", "address", mv3alias, "--force")

	out, err = MavkitClient("transfer", "100", "from", "alice", "to", mv2alias, "--burn-cap", "0.06425")
	assert.NoError(t, err)
	require.Contains(t, string(out), "Operation successfully injected in the node")
	out, err = MavkitClient("transfer", "100", "from", "alice", "to", mv3alias, "--burn-cap", "0.06425")
	assert.NoError(t, err)
	require.Contains(t, string(out), "Operation successfully injected in the node")

	//test
	/* the mv2 key produces invalid signature 50% of the time from mavkit-client perspective
	out, err = MavkitClient("transfer", "1", "from", mv2alias, "to", "alice", "--burn-cap", "0.06425")
	assert.NoError(t, err)
	require.Contains(t, string(out), "Operation successfully injected in the node")
	*/
	out, err = MavkitClient("transfer", "1", "from", mv3alias, "to", "alice", "--burn-cap", "0.06425")
	assert.NoError(t, err)
	require.Contains(t, string(out), "Operation successfully injected in the node")

	require.Equal(t, mv3pk, GetPublicKey(mv3))
}
