package node_test

import (
	"os"
	"testing"
	"time"

	"github.com/filecoin-project/lotus/chain/actors/builtin/miner"

	builder "github.com/filecoin-project/lotus/node/test"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/lotus/lib/lotuslog"
	miner0 "github.com/filecoin-project/specs-actors/actors/builtin/miner"
	power0 "github.com/filecoin-project/specs-actors/actors/builtin/power"
	verifreg0 "github.com/filecoin-project/specs-actors/actors/builtin/verifreg"
	logging "github.com/ipfs/go-log/v2"

	"github.com/filecoin-project/lotus/api/test"
)

func init() {
	_ = logging.SetLogLevel("*", "INFO")

	power0.ConsensusMinerMinPower = big.NewInt(2048)
	miner0.SupportedProofTypes = map[abi.RegisteredSealProof]struct{}{
		abi.RegisteredSealProof_StackedDrg2KiBV1: {},
	}
	verifreg0.MinVerifiedDealSize = big.NewInt(256)
}

func TestAPI(t *testing.T) {
	test.TestApis(t, builder.Builder)
}

func TestAPIRPC(t *testing.T) {
	test.TestApis(t, builder.RPCBuilder)
}

func TestAPIDealFlow(t *testing.T) {
	logging.SetLogLevel("miner", "ERROR")
	logging.SetLogLevel("chainstore", "ERROR")
	logging.SetLogLevel("chain", "ERROR")
	logging.SetLogLevel("sub", "ERROR")
	logging.SetLogLevel("storageminer", "ERROR")

	t.Run("TestDealFlow", func(t *testing.T) {
		test.TestDealFlow(t, builder.MockSbBuilder, 10*time.Millisecond, false, false)
	})
	t.Run("WithExportedCAR", func(t *testing.T) {
		test.TestDealFlow(t, builder.MockSbBuilder, 10*time.Millisecond, true, false)
	})
	t.Run("TestDoubleDealFlow", func(t *testing.T) {
		test.TestDoubleDealFlow(t, builder.MockSbBuilder, 10*time.Millisecond)
	})
	t.Run("TestFastRetrievalDealFlow", func(t *testing.T) {
		test.TestFastRetrievalDealFlow(t, builder.MockSbBuilder, 10*time.Millisecond)
	})
}

func TestAPIDealFlowReal(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	lotuslog.SetupLogLevels()
	logging.SetLogLevel("miner", "ERROR")
	logging.SetLogLevel("chainstore", "ERROR")
	logging.SetLogLevel("chain", "ERROR")
	logging.SetLogLevel("sub", "ERROR")
	logging.SetLogLevel("storageminer", "ERROR")

	// TODO: Do this better.
	miner.PreCommitChallengeDelay = 5
	miner0.PreCommitChallengeDelay = 5

	t.Run("basic", func(t *testing.T) {
		test.TestDealFlow(t, builder.Builder, time.Second, false, false)
	})

	t.Run("fast-retrieval", func(t *testing.T) {
		test.TestDealFlow(t, builder.Builder, time.Second, false, true)
	})

	t.Run("retrieval-second", func(t *testing.T) {
		test.TestSenondDealRetrieval(t, builder.Builder, time.Second)
	})
}

func TestDealMining(t *testing.T) {
	logging.SetLogLevel("miner", "ERROR")
	logging.SetLogLevel("chainstore", "ERROR")
	logging.SetLogLevel("chain", "ERROR")
	logging.SetLogLevel("sub", "ERROR")
	logging.SetLogLevel("storageminer", "ERROR")

	test.TestDealMining(t, builder.MockSbBuilder, 50*time.Millisecond, false)
}

func TestPledgeSectors(t *testing.T) {
	logging.SetLogLevel("miner", "ERROR")
	logging.SetLogLevel("chainstore", "ERROR")
	logging.SetLogLevel("chain", "ERROR")
	logging.SetLogLevel("sub", "ERROR")
	logging.SetLogLevel("storageminer", "ERROR")

	t.Run("1", func(t *testing.T) {
		test.TestPledgeSector(t, builder.MockSbBuilder, 50*time.Millisecond, 1)
	})

	t.Run("100", func(t *testing.T) {
		test.TestPledgeSector(t, builder.MockSbBuilder, 50*time.Millisecond, 100)
	})

	t.Run("1000", func(t *testing.T) {
		if testing.Short() { // takes ~16s
			t.Skip("skipping test in short mode")
		}

		test.TestPledgeSector(t, builder.MockSbBuilder, 50*time.Millisecond, 1000)
	})
}

func TestWindowedPost(t *testing.T) {
	if os.Getenv("LOTUS_TEST_WINDOW_POST") != "1" {
		t.Skip("this takes a few minutes, set LOTUS_TEST_WINDOW_POST=1 to run")
	}

	logging.SetLogLevel("miner", "ERROR")
	logging.SetLogLevel("chainstore", "ERROR")
	logging.SetLogLevel("chain", "ERROR")
	logging.SetLogLevel("sub", "ERROR")
	logging.SetLogLevel("storageminer", "ERROR")

	test.TestWindowPost(t, builder.MockSbBuilder, 2*time.Millisecond, 10)
}

func TestCCUpgrade(t *testing.T) {
	logging.SetLogLevel("miner", "ERROR")
	logging.SetLogLevel("chainstore", "ERROR")
	logging.SetLogLevel("chain", "ERROR")
	logging.SetLogLevel("sub", "ERROR")
	logging.SetLogLevel("storageminer", "ERROR")

	test.TestCCUpgrade(t, builder.MockSbBuilder, 5*time.Millisecond)
}

func TestPaymentChannels(t *testing.T) {
	logging.SetLogLevel("miner", "ERROR")
	logging.SetLogLevel("chainstore", "ERROR")
	logging.SetLogLevel("chain", "ERROR")
	logging.SetLogLevel("sub", "ERROR")
	logging.SetLogLevel("pubsub", "ERROR")
	logging.SetLogLevel("storageminer", "ERROR")

	test.TestPaymentChannels(t, builder.MockSbBuilder, 5*time.Millisecond)
}
