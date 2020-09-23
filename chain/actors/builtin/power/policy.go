package power

import (
	"github.com/filecoin-project/go-state-types/abi"

	power0 "github.com/filecoin-project/specs-actors/actors/builtin/power"
)

// SetConsensusMinerMinPower sets the minimum power of an individual miner must
// meet for leader election, across all actor versions. This should only be used
// for testing.
func SetConsensusMinerMinPower(p abi.StoragePower) {
	power0.ConsensusMinerMinPower = p
}
