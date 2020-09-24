package simpleretr

import (
	//"context"

	"github.com/libp2p/go-libp2p-core/network"
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
	HandleStream(stream network.Stream)
}
