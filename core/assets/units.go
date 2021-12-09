package assets

import (
	"math/big"

	"github.com/klaytn/klaytn/params"
)

func Wei(n int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(n), big.NewInt(params.Peb))
}

func GWei(n int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(n), big.NewInt(params.Ston))
}

func Ether(n int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(n), big.NewInt(params.KLAY))
}
