package integrationtest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

//there are enough existing integration tests using mv1 and mv3 addresses that it would be redundant to do so here

func TestTz4(t *testing.T) {
	//flexmasa does not start when we try to use a mv4 bootstrap account, so, we have to fund it first
	out, err := MavkitClient("-w", "1", "transfer", "200", "from", "alice", "to", "mv4alias", "--burn-cap", "0.06425")
	require.NoError(t, err)
	require.Contains(t, string(out), "Operation successfully injected in the node")

	out, err = MavkitClient("-w", "1", "transfer", "100", "from", "mv4alias", "to", "alice", "--burn-cap", "0.06425")
	require.NoError(t, err)
	require.Contains(t, string(out), "Operation successfully injected in the node")
}

func TestTz2(t *testing.T) {
	out, err := MavkitClient("-w", "1", "transfer", "100", "from", "mv2alias", "to", "alice", "--burn-cap", "0.06425")
	require.NoError(t, err)
	require.Contains(t, string(out), "Operation successfully injected in the node")
}
