# exchain-ethereum-compatible

## Background
ExChain use a different method to calculate the hash value of a transaction. Instead of rlp encode and keccak256 in ethereum format, ExChain imports a new object encoding specification, called [go-amino](https://github.com/tendermint/go-amino), and performs the SHA256 algorithm on the amino-encoded data to calculate the real hash.

Because of [exchain](github.com/okex/exchain) importing an outdated go-ethereum, some projects relying on a higher version of go-ethereum who want to integrate encoded functions might encounter problems of dependencies corrupting in go module

## Solution
To be compatible with hash between ethereum and exchain, and make less effort to migrate the project from ethereum to exchain, this package is used for calculate the real hash of an evm transaction.

### usage
go.mod in your project
```go
require github.com/smartcontractkit/chainlink/core/external/exchain-celo-compatible v1.0.2
```

Instead of signtx.Hash(), using utils.Hash(signtx)

```go
import (
    "github.com/celo-org/celo-blockchain/core/types"
    "github.com/smartcontractkit/chainlink/core/external/exchain-celo-compatible/utils"
)

func Test() {
   //...
	
    unsignedTx := types.NewTransaction(	nonce, common.HexToAddress("0x79BE5cc37B7e17594028BbF5d43875FDbed417da"), big.NewInt(1e18), uint64(3000000), gasPrice, nil);
    signedTx, err := types.SignTx(unsignedTx, types.NewLondonSigner(chainID), privateKey)

    //...
    
    // signedTx.Hash() works on ethereum, but not on ExChain
    // So try to use utils.Hash(xxx)
    hash, _ := utils.Hash(signedTx)
    fmt.Println(hash.ToString())
}

```
