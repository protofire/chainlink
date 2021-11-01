package keeper

import (
	ethereum "github.com/celo-org/celo-blockchain"
)

func (rs *RegistrySynchronizer) ExportedFullSync() {
	rs.fullSync()
}

func (rs *RegistrySynchronizer) ExportedProcessLogs() {
	rs.processLogs()
}

func (ex *UpkeepExecuter) ExportedConstructCheckUpkeepCallMsg(upkeep UpkeepRegistration) (ethereum.CallMsg, error) {
	return ex.constructCheckUpkeepCallMsg(upkeep)
}
