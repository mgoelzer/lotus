package verifreg

import (
	"github.com/filecoin-project/go-state-types/abi"

	verifreg0 "github.com/filecoin-project/specs-actors/actors/builtin/verifreg"
)

// SetMinVerifiedDealSize sets the minimum size of a verified deal. This should
// only be used for testing.
func SetMinVerifiedDealSize(size abi.StoragePower) {
	verifreg0.MinVerifiedDealSize = size
}
