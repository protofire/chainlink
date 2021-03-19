package log

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type (
	ethSubscriber struct {
		ethClient eth.Client
		config    Config
		chStop    chan struct{}
		chDone    chan struct{}
	}
)

func newEthSubscriber(ethClient eth.Client, config Config) *ethSubscriber {
	return &ethSubscriber{
		ethClient: ethClient,
		config:    config,
		chStop:    make(chan struct{}),
		chDone:    make(chan struct{}),
	}
}
func (b *ethSubscriber) Stop() error {
	close(b.chStop)
	<-b.chDone
	return nil
}

func (b *ethSubscriber) backfillLogs(addresses []common.Address, topics []common.Hash) (chBackfilledLogs chan types.Log, abort bool) {
	if len(addresses) == 0 {
		ch := make(chan types.Log)
		close(ch)
		return ch, false
	}

	ctx, cancel := utils.ContextFromChan(b.chStop)
	defer cancel()

	utils.RetryWithBackoff(ctx, func() (retry bool) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		latestBlock, err := b.ethClient.HeaderByNumber(ctx, nil)
		if err != nil {
			logger.Errorw("Log subscriber backfill: could not fetch latest block header", "error", err)
			return true
		} else if latestBlock == nil {
			logger.Warn("got nil block header")
			return true
		}
		currentHeight := uint64(latestBlock.Number)

		// Backfill from `backfillDepth` blocks ago.  It's up to the subscribers to
		// filter out logs they've already dealt with.
		fromBlock := currentHeight - b.config.BlockBackfillDepth()
		if fromBlock > currentHeight {
			fromBlock = 0 // Overflow protection
		}

		q := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(fromBlock)),
			Addresses: addresses,
			Topics:    [][]common.Hash{topics},
		}

		logs, err := b.ethClient.FilterLogs(ctx, q)
		if err != nil {
			logger.Errorw("Log subscriber backfill: could not fetch logs", "error", err)
			return true
		}

		chBackfilledLogs = make(chan types.Log)
		go func() {
			defer close(chBackfilledLogs)
			for _, log := range logs {
				select {
				case chBackfilledLogs <- log:
				case <-b.chStop:
					return
				}
			}
		}()

		return false
	})
	select {
	case <-b.chStop:
		abort = true
	default:
		abort = false
	}
	return
}

// createSubscription creates a new log subscription starting at the current block.  If previous logs
// are needed, they must be obtained through backfilling, as subscriptions can only be started from
// the current head.
func (b *ethSubscriber) createSubscription(addresses []common.Address, topics []common.Hash) (sub managedSubscription, abort bool) {
	if len(addresses) == 0 {
		return newNoopSubscription(), false
	}

	ctx, cancel := utils.ContextFromChan(b.chStop)
	defer cancel()

	utils.RetryWithBackoff(ctx, func() (retry bool) {

		filterQuery := ethereum.FilterQuery{
			Addresses: addresses,
			Topics:    [][]common.Hash{topics},
		}
		chRawLogs := make(chan types.Log)

		ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()

		innerSub, err := b.ethClient.SubscribeFilterLogs(ctx, filterQuery, chRawLogs)
		if err != nil {
			logger.Errorw("Log subscriber could not create subscription to Ethereum node", "error", err)
			return true
		}

		sub = managedSubscriptionImpl{
			subscription: innerSub,
			chRawLogs:    chRawLogs,
		}
		return false
	})
	select {
	case <-b.chStop:
		abort = true
	default:
		abort = false
	}
	return
}

// A managedSubscription acts as wrapper for the Subscription. Specifically, the
// managedSubscription closes the log channel as soon as the unsubscribe request is made
type managedSubscription interface {
	Err() <-chan error
	Logs() chan types.Log
	Unsubscribe()
}

type managedSubscriptionImpl struct {
	subscription ethereum.Subscription
	chRawLogs    chan types.Log
}

func (sub managedSubscriptionImpl) Err() <-chan error {
	return sub.subscription.Err()
}

func (sub managedSubscriptionImpl) Logs() chan types.Log {
	return sub.chRawLogs
}

func (sub managedSubscriptionImpl) Unsubscribe() {
	sub.subscription.Unsubscribe()
	close(sub.chRawLogs)
}

type noopSubscription struct {
	chRawLogs chan types.Log
}

func newNoopSubscription() noopSubscription {
	return noopSubscription{make(chan types.Log)}
}

func (b noopSubscription) Err() <-chan error    { return nil }
func (b noopSubscription) Logs() chan types.Log { return b.chRawLogs }
func (b noopSubscription) Unsubscribe()         { close(b.chRawLogs) }

// ListenerJobID returns the appropriate job ID for a listener
func ListenerJobID(listener Listener) interface{} {
	if listener.IsV2Job() {
		return listener.JobIDV2()
	}
	return listener.JobID()
}
