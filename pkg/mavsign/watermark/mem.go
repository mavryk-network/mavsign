package watermark

import (
	"context"
	"sync"

	mv "github.com/mavryk-network/mavbingo/v2"
	"github.com/mavryk-network/mavbingo/v2/crypt"
	"github.com/mavryk-network/mavbingo/v2/protocol"
	"github.com/mavryk-network/mavsign/pkg/config"
	"github.com/mavryk-network/mavsign/pkg/mavsign/request"
	"gopkg.in/yaml.v3"
)

// InMemory keep previous operation in memory
type InMemory struct {
	chains map[mv.ChainID]delegateMap
	mtx    sync.Mutex
}

// IsSafeToSign return true if this msgID is safe to sign
func (w *InMemory) IsSafeToSign(ctx context.Context, pkh crypt.PublicKeyHash, req protocol.SignRequest, digest *crypt.Digest) error {
	w.mtx.Lock()
	defer w.mtx.Unlock()
	return w.isSafeToSignUnlocked(pkh, req, digest)
}

func (w *InMemory) isSafeToSignUnlocked(pkh crypt.PublicKeyHash, req protocol.SignRequest, digest *crypt.Digest) error {
	m, ok := req.(request.WithWatermark)
	if !ok {
		// watermark is not required
		return nil
	}

	if w.chains == nil {
		w.chains = make(map[mv.ChainID]delegateMap)
	}

	delegates, ok := w.chains[*m.GetChainID()]
	if !ok {
		delegates = make(delegateMap)
		w.chains[*m.GetChainID()] = delegates
	}

	requests, ok := delegates.Get(pkh)
	if !ok {
		requests = make(requestMap)
		delegates.Insert(pkh, requests)
	}

	watermark := request.NewWatermark(m, digest)
	if stored, ok := requests[req.SignRequestKind()]; ok {
		if !watermark.Validate(stored) {
			return ErrWatermark
		}
	}
	requests[m.SignRequestKind()] = watermark
	return nil
}

func init() {
	RegisterWatermark("mem", func(context.Context, *yaml.Node, *config.Config) (Watermark, error) {
		return new(InMemory), nil
	})
}
