package simpleretr

import (
	//"context"

	inet "github.com/libp2p/go-libp2p-core/network"
	//"github.com/libp2p/go-libp2p-core/peer"
)

// Server is the responder side of the /fil/simple-retrieve. It accepts
// requests from the client and replies with the corresponding response
// struct.
type Server interface {
	// HandleStream is the protocol handler to be registered on a libp2p
	// protocol router.
	//
	// In the current version of the protocol, streams are single-use. The
	// server will read a single Request, and will respond with a single
	// Response. It will dispose of the stream straight after.
	HandleStream(stream inet.Stream)
}

// Client is the requesting side of the ChainExchange protocol. It acts as
// a proxy for other components to request chain data from peers. It is chiefly
// used by the Syncer.
type Client interface {

	// TODO:  delete this
	//	// GetBlocks fetches block headers from the network, from the provided
	//	// tipset *backwards*, returning as many tipsets as the count parameter,
	//	// or less.
	//	GetBlocks(ctx context.Context, tsk types.TipSetKey, count int) ([]*types.TipSet, error)

}
