//go:build !integration

package watermark

import (
	"context"
	"fmt"
	"os"
	"testing"

	mv "github.com/mavryk-network/mavbingo/v2"
	"github.com/mavryk-network/mavbingo/v2/crypt"
	"github.com/mavryk-network/mavbingo/v2/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dummyMsg struct {
	kind  string
	level int32
	round int32
}

func (r *dummyMsg) SignRequestKind() string { return r.kind }
func (r *dummyMsg) GetChainID() *mv.ChainID { return &mv.ChainID{} }
func (r *dummyMsg) GetLevel() int32         { return r.level }
func (r *dummyMsg) GetRound() int32         { return r.round }

type testCase struct {
	pkh       crypt.PublicKeyHash
	req       protocol.SignRequest
	reqDigest crypt.Digest
	expectErr bool
}

func TestWatermark(t *testing.T) {
	cases := []testCase{
		{
			pkh: &mv.Ed25519PublicKeyHash{0},
			req: &dummyMsg{
				kind:  "kind0",
				level: 124,
			},
			reqDigest: crypt.Digest{0},
			expectErr: false,
		},
		{
			pkh: &mv.Ed25519PublicKeyHash{0},
			req: &dummyMsg{
				kind:  "kind0",
				level: 124,
			},
			reqDigest: crypt.Digest{1},
			expectErr: true, // same level
		},
		{
			pkh: &mv.Ed25519PublicKeyHash{0},
			req: &dummyMsg{
				kind:  "kind0",
				level: 123,
			},
			reqDigest: crypt.Digest{2},
			expectErr: true, // level below
		},
		{
			pkh: &mv.Ed25519PublicKeyHash{0},
			req: &dummyMsg{
				kind:  "kind0",
				level: 124,
			},
			reqDigest: crypt.Digest{0},
			expectErr: false, // repeated request
		},
		{
			pkh: &mv.Ed25519PublicKeyHash{1},
			req: &dummyMsg{
				kind:  "kind0",
				level: 124,
			},
			reqDigest: crypt.Digest{3},
			expectErr: false, // different delegate
		},
		{
			pkh: &mv.Ed25519PublicKeyHash{1},
			req: &dummyMsg{
				kind:  "kind0",
				level: 125,
			},
			reqDigest: crypt.Digest{4},
			expectErr: false,
		},
		{
			pkh: &mv.Ed25519PublicKeyHash{0},
			req: &dummyMsg{
				kind:  "kind1",
				level: 124,
			},
			reqDigest: crypt.Digest{0},
			expectErr: false, // different kind
		},
	}

	t.Run("memory", func(t *testing.T) {
		var wm InMemory
		for i, c := range cases {
			t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
				err := wm.IsSafeToSign(context.Background(), c.pkh, c.req, &c.reqDigest)
				if c.expectErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("file", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "watermark")
		require.NoError(t, err)
		wm, err := NewFileWatermark(dir)
		require.NoError(t, err)
		for i, c := range cases {
			t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
				err := wm.IsSafeToSign(context.Background(), c.pkh, c.req, &c.reqDigest)
				if c.expectErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}
