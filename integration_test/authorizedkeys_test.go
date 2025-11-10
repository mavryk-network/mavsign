package integrationtest

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthorizedKeys(t *testing.T) {
	var c Config
	c.Read()
	c.Server.Keys = []string{"edpkujLb5ZCZ2gprnRzE9aVHKZfx9A8EtWu2xxkwYSjBUJbesJ9rWE"}
	backup_then_update_config(c)
	defer restore_config()
	restart_mavsign()

	out, err := MavkitClient("-w", "1", "transfer", "1", "from", "alice", "to", "bob")
	require.NotNil(t, err)
	require.Contains(t, string(out), "remote signer expects authentication signature, but no authorized key was found in the wallet")

	out, err = MavkitClient("import", "secret", "key", "auth", "unencrypted:edsk3ZAm9nqEo7qNugo2wcmxWnbDe7oUUmHt5UJYDdqwucsaHTAsVQ", "--force")
	defer MavkitClient("forget", "address", "auth", "--force")
	if err != nil {
		log.Println("failed to import auth key: " + err.Error() + string(out))
	}
	require.Nil(t, err)

	out, err = MavkitClient("-w", "1", "transfer", "1", "from", "alice", "to", "bob")
	require.Nil(t, err)
	require.Contains(t, string(out), "Operation successfully injected in the node.")
}
