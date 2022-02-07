package client

import (
	"context"
	"math/big"

	celo "github.com/celo-org/celo-blockchain"
	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/core/types"
	"github.com/celo-org/celo-blockchain/rpc"
	"github.com/pkg/errors"
)

var _ Node = (*erroringNode)(nil)

type erroringNode struct {
	errMsg string
}

func (e *erroringNode) ChainID(ctx context.Context) (chainID *big.Int, err error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) Dial(ctx context.Context) error {
	return errors.New(e.errMsg)
}

func (e *erroringNode) Close() {}

func (e *erroringNode) Verify(ctx context.Context, expectedChainID *big.Int) (err error) {
	return errors.New(e.errMsg)
}

func (e *erroringNode) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return errors.New(e.errMsg)
}

func (e *erroringNode) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	return errors.New(e.errMsg)
}

func (e *erroringNode) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return errors.New(e.errMsg)
}

func (e *erroringNode) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return 0, errors.New(e.errMsg)
}

func (e *erroringNode) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	return 0, errors.New(e.errMsg)
}

func (e *erroringNode) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) FilterLogs(ctx context.Context, q celo.FilterQuery) ([]types.Log, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) SubscribeFilterLogs(ctx context.Context, q celo.FilterQuery, ch chan<- types.Log) (celo.Subscription, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) EstimateGas(ctx context.Context, call celo.CallMsg) (uint64, error) {
	return 0, errors.New(e.errMsg)
}

func (e *erroringNode) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) CallContract(ctx context.Context, msg celo.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) HeaderByNumber(_ context.Context, _ *big.Int) (*types.Header, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) EthSubscribe(ctx context.Context, channel interface{}, args ...interface{}) (celo.Subscription, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) String() string {
	return "<erroring node>"
}

func (e *erroringNode) State() NodeState {
	return NodeStateDead
}
