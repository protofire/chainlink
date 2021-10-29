package internal

import "github.com/tendermint/go-amino"

var moduleCdc *amino.Codec

func init() {
	cdc := amino.NewCodec()
	cdc.RegisterConcrete(MsgEthereumTx{}, "ethermint/MsgEthereumTx", nil)
	cdc.RegisterConcrete(TxData{}, "ethermint/TxData", nil)
	cdc.Seal()
	moduleCdc = cdc
}

func GetModuleCdc() *amino.Codec {
	return moduleCdc
}
