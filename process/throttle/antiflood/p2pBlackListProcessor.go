package antiflood

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-go/storage"
)

type p2pBlackListProcessor struct {
	cacher                     storage.Cacher
	blacklistHandler           process.BlackListHandler
	thresholdNumReceivedFlood  uint32
	thresholdSizeReceivedFlood uint64
	numFloodingRounds          uint32
}

// NewP2pBlackListProcessor creates a new instance of p2pQuotaBlacklistProcessor able to determine
// a flooding peer and mark it accordingly
func NewP2pBlackListProcessor(
	cacher storage.Cacher,
	blacklistHandler process.BlackListHandler,
	thresholdNumReceivedFlood uint32,
	thresholdSizeReceivedFlood uint64,
	numFloodingRounds uint32,
) (*p2pBlackListProcessor, error) {

	if check.IfNil(cacher) {
		return nil, fmt.Errorf("%w, NewP2pBlackListProcessor", process.ErrNilCacher)
	}
	if check.IfNil(blacklistHandler) {
		return nil, fmt.Errorf("%w, NewP2pBlackListProcessor", process.ErrNilBlackListHandler)
	}
	if thresholdNumReceivedFlood == 0 {
		return nil, fmt.Errorf("%w, thresholdNumReceivedFlood == 0", process.ErrInvalidValue)
	}
	if thresholdSizeReceivedFlood == 0 {
		return nil, fmt.Errorf("%w, thresholdSizeReceivedFlood == 0", process.ErrInvalidValue)
	}
	if numFloodingRounds == 0 {
		return nil, fmt.Errorf("%w, numFloodingRounds == 0", process.ErrInvalidValue)
	}

	return &p2pBlackListProcessor{
		cacher:                     cacher,
		blacklistHandler:           blacklistHandler,
		thresholdNumReceivedFlood:  thresholdNumReceivedFlood,
		thresholdSizeReceivedFlood: thresholdSizeReceivedFlood,
		numFloodingRounds:          numFloodingRounds,
	}, nil
}

// ResetStatistics checks if an identifier reached its maximum flooding rounds. If it did, it will remove its
// cached information and adds it to the black list handler
func (pbp *p2pBlackListProcessor) ResetStatistics() {
	keys := pbp.cacher.Keys()
	for _, key := range keys {
		obj, ok := pbp.cacher.Peek(key)
		if !ok {
			pbp.cacher.Remove(key)
			continue
		}

		val, ok := obj.(uint32)
		if !ok {
			pbp.cacher.Remove(key)
			continue
		}

		if val >= pbp.numFloodingRounds {
			pbp.cacher.Remove(key)
			_ = pbp.blacklistHandler.Add(string(key))
		}
	}
}

// AddQuota checks if the received quota for an identifier has exceeded the set thresholds
func (pbp *p2pBlackListProcessor) AddQuota(identifier string, numReceived uint32, sizeReceived uint64, _ uint32, _ uint64) {
	isFloodingPeer := numReceived >= pbp.thresholdNumReceivedFlood || sizeReceived >= pbp.thresholdSizeReceivedFlood
	if isFloodingPeer {
		pbp.incrementStatsFloodingPeer(identifier)
	}
}

func (pbp *p2pBlackListProcessor) incrementStatsFloodingPeer(identifier string) {
	obj, ok := pbp.cacher.Get([]byte(identifier))
	if !ok {
		pbp.cacher.Put([]byte(identifier), uint32(1))
		return
	}

	val, ok := obj.(uint32)
	if !ok {
		pbp.cacher.Put([]byte(identifier), uint32(1))
		return
	}

	pbp.cacher.Put([]byte(identifier), val+1)
}

// SetGlobalQuota does nothing (here to comply with QuotaStatusHandler interface)
func (pbp *p2pBlackListProcessor) SetGlobalQuota(_ uint32, _ uint64, _ uint32, _ uint64) {}

// IsInterfaceNil returns true if there is no value under the interface
func (pbp *p2pBlackListProcessor) IsInterfaceNil() bool {
	return pbp == nil
}