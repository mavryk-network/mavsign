package request

import (
	mv "github.com/mavryk-network/mavbingo/v2"
	"github.com/mavryk-network/mavbingo/v2/crypt"
	"github.com/mavryk-network/mavbingo/v2/protocol"
)

type WithWatermark interface {
	protocol.SignRequest
	GetChainID() *mv.ChainID
	GetLevel() int32
	GetRound() int32
}

type Watermark struct {
	Level int32                          `json:"level"`
	Round int32                          `json:"round"`
	Hash  mv.Option[mv.BlockPayloadHash] `json:"hash"`
}

func NewWatermark(req WithWatermark, hash *crypt.Digest) *Watermark {
	return &Watermark{
		Level: req.GetLevel(),
		Round: req.GetRound(),
		Hash:  mv.Some((mv.BlockPayloadHash)(*hash)),
	}
}

func (l *Watermark) Validate(stored *Watermark) bool {
	if l.Hash.IsSome() && stored.Hash.IsSome() && l.Hash.Unwrap() == stored.Hash.Unwrap() {
		return true
	}
	var diff int32
	if d := l.Level - stored.Level; d == 0 {
		diff = l.Round - stored.Round
	} else {
		diff = d
	}
	return diff > 0
}

var (
	_ WithWatermark = (*protocol.BlockSignRequest)(nil)
	_ WithWatermark = (*protocol.PreendorsementSignRequest)(nil)
	_ WithWatermark = (*protocol.EndorsementSignRequest)(nil)
)
