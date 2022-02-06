package handler

import (
	"context"
	"log"

	"github.com/celo-org/celo-blockchain/accounts/abi/bind"
	"github.com/celo-org/celo-blockchain/core/types"
)

func waitDeployment(ctx context.Context, client bind.DeployBackend, tx *types.Transaction) {
	if _, err := bind.WaitDeployed(ctx, client, tx); err != nil {
		log.Fatal("WaitDeployed failed: ", err)
	}
}

func waitTx(ctx context.Context, client bind.DeployBackend, tx *types.Transaction) {
	if _, err := bind.WaitMined(ctx, client, tx); err != nil {
		log.Fatal("WaitDeployed failed: ", err)
	}
}
