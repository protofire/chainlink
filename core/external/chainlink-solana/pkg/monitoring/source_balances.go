package monitoring

import (
	"context"
	"fmt"
	"sync"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	relayMonitoring "github.com/smartcontractkit/chainlink-relay/pkg/monitoring"
	pkgSolana "github.com/smartcontractkit/chainlink/core/external/chainlink-solana/pkg/solana"
)

func NewBalancesSourceFactory(
	solanaConfig SolanaConfig,
	log relayMonitoring.Logger,
) relayMonitoring.SourceFactory {
	client := rpc.New(solanaConfig.RPCEndpoint)
	return &balancesSourceFactory{
		client,
		log,
	}
}

type balancesSourceFactory struct {
	client *rpc.Client
	log    relayMonitoring.Logger
}

func (s *balancesSourceFactory) NewSource(
	chainConfig relayMonitoring.ChainConfig,
	feedConfig relayMonitoring.FeedConfig,
) (relayMonitoring.Source, error) {
	solanaConfig, ok := chainConfig.(SolanaConfig)
	if !ok {
		return nil, fmt.Errorf("expected chainConfig to be of type SolanaConfig not %T", chainConfig)
	}
	solanaFeedConfig, ok := feedConfig.(SolanaFeedConfig)
	if !ok {
		return nil, fmt.Errorf("expected feedConfig to be of type SolanaFeedConfig not %T", feedConfig)
	}
	return &balancesSource{
		s.client,
		s.log,
		solanaConfig,
		solanaFeedConfig,
	}, nil
}

type balancesSource struct {
	client       *rpc.Client
	log          relayMonitoring.Logger
	solanaConfig SolanaConfig
	feedConfig   SolanaFeedConfig
}

type Balances struct {
	Values    map[string]uint64
	Addresses map[string]solana.PublicKey
}

func (s *balancesSource) Fetch(ctx context.Context) (interface{}, error) {
	state, _, err := pkgSolana.GetState(ctx, s.client, s.feedConfig.StateAccount, rpc.CommitmentConfirmed)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract state: %w", err)
	}
	isErr := false
	balances := Balances{
		Values:    make(map[string]uint64),
		Addresses: make(map[string]solana.PublicKey),
	}
	balancesMu := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	wg.Add(len(BalanceAccountNames))
	for key, address := range map[string]solana.PublicKey{
		"contract":                    s.feedConfig.ContractAddress,
		"state":                       s.feedConfig.StateAccount,
		"transmissions":               state.Transmissions,
		"token_vault":                 state.Config.TokenVault,
		"requester_access_controller": state.Config.RequesterAccessController,
		"billing_access_controller":   state.Config.BillingAccessController,
	} {
		go func(key string, address solana.PublicKey) {
			defer wg.Done()
			res, err := s.client.GetBalance(ctx, address, rpc.CommitmentProcessed)
			balancesMu.Lock()
			defer balancesMu.Unlock()
			if err != nil {
				s.log.Errorw("failed to read the sol balance", "key", key, "address", address.String(), "error", err)
				isErr = true
				return
			}
			balances.Values[key] = res.Value
			balances.Addresses[key] = address
		}(key, address)
	}

	wg.Wait()
	if isErr {
		return Balances{}, fmt.Errorf("error while fetching balances")
	}
	return balances, nil
}
