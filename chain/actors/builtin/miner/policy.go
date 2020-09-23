package miner

import (
	"github.com/filecoin-project/go-state-types/abi"

	miner0 "github.com/filecoin-project/specs-actors/actors/builtin/miner"
)

// SetSupportedProofTypes sets supported proof types, across all actor versions.
// This should only be used for testing.
func SetSupportedProofTypes(types ...abi.RegisteredSealProof) {
	newTypes := make(map[abi.RegisteredSealProof]struct{}, len(types))
	for _, t := range types {
		newTypes[t] = struct{}{}
	}
	// Set for all miner versions.
	miner0.SupportedProofTypes = newTypes
}

// AddSupportedProofTypes sets supported proof types, across all actor versions.
// This should only be used for testing.
func AddSupportedProofTypes(types ...abi.RegisteredSealProof) {
	newTypes := make(map[abi.RegisteredSealProof]struct{}, len(types))
	for _, t := range types {
		newTypes[t] = struct{}{}
	}
	// Set for all miner versions.
	miner0.SupportedProofTypes = newTypes
}

// SetPreCommitChallengeDelay sets the pre-commit challenge delay across all
// actors versions. Use for testing.
func SetPreCommitChallengeDelay(delay abi.ChainEpoch) {
	// Set for all miner versions.
	miner0.PreCommitChallengeDelay = delay
	PreCommitChallengeDelay = delay
}
