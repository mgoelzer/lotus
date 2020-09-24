package simpleretr

import (
	//"bufio"
	//"context"
	//"fmt"
	//"math/rand"
	//"time"

	"github.com/libp2p/go-libp2p-core/host"
	//"github.com/libp2p/go-libp2p-core/network"
	//"github.com/libp2p/go-libp2p-core/peer"

	//"go.opencensus.io/trace"
	//"go.uber.org/fx"
	//"golang.org/x/xerrors"

	//cborutil "github.com/filecoin-project/go-cbor-util"

	//"github.com/filecoin-project/lotus/build"
)

// client implements the interface simpleretr.Client, which is the initiator 
// and data receiver side of the /fil/simple-retrieve protocol.
type client struct {
	// client's libp2p host (used to contact the server peer)
	host host.Host

	// multiaddr of the server we wish to interact with
	serverMultiaddr string
}

// TODO:  what is the purpose of this?
var _ Client = (*client)(nil)

func NewClient(host host.Host, serverMultiaddr string) Client {
	return &client{
		host: host,
		serverMultiaddr: serverMultiaddr,
	}
}

// Sends a request to the server, and returns the server's response. The request
// is a typed request struct that will be serialized to a JSON object on the wrire, 
// while the server's response is deserialized from JSON here and returned as a 
// typed response struct.
/*func (c *client) sendRequestToPeer(ctx context.Context, peer peer.ID, req *Request) (_ *Response, err error) {
	// Trace code.
	ctx, span := trace.StartSpan(ctx, "sendRequestToPeer")
	defer span.End()
	if span.IsRecordingEvents() {
		span.AddAttributes(
			trace.StringAttribute("peer", peer.Pretty()),
		)
	}
	defer func() {
		if err != nil {
			if span.IsRecordingEvents() {
				span.SetStatus(trace.Status{
					Code:    5,
					Message: err.Error(),
				})
			}
		}
	}()
	// -- TRACE --

	supported, err := c.host.Peerstore().SupportsProtocols(peer, BlockSyncProtocolID, ChainExchangeProtocolID)
	if err != nil {
		c.RemovePeer(peer)
		return nil, xerrors.Errorf("failed to get protocols for peer: %w", err)
	}
	if len(supported) == 0 || (supported[0] != BlockSyncProtocolID && supported[0] != ChainExchangeProtocolID) {
		return nil, xerrors.Errorf("peer %s does not support protocols %s",
			peer, []string{BlockSyncProtocolID, ChainExchangeProtocolID})
	}

	connectionStart := build.Clock.Now()

	// Open stream to peer.
	stream, err := c.host.NewStream(
		network.WithNoDial(ctx, "should already have connection"),
		peer,
		ChainExchangeProtocolID, BlockSyncProtocolID)
	if err != nil {
		c.RemovePeer(peer)
		return nil, xerrors.Errorf("failed to open stream to peer: %w", err)
	}

	// Write request.
	_ = stream.SetWriteDeadline(time.Now().Add(WriteReqDeadline))
	if err := cborutil.WriteCborRPC(stream, req); err != nil {
		_ = stream.SetWriteDeadline(time.Time{})
		c.peerTracker.logFailure(peer, build.Clock.Since(connectionStart), req.Length)
		// FIXME: Should we also remove peer here?
		return nil, err
	}
	_ = stream.SetWriteDeadline(time.Time{}) // clear deadline // FIXME: Needs
	//  its own API (https://github.com/libp2p/go-libp2p-core/issues/162).

	// Read response.
	var res Response
	err = cborutil.ReadCborRPC(
		bufio.NewReader(incrt.New(stream, ReadResMinSpeed, ReadResDeadline)),
		&res)
	if err != nil {
		c.peerTracker.logFailure(peer, build.Clock.Since(connectionStart), req.Length)
		return nil, xerrors.Errorf("failed to read chainxchg response: %w", err)
	}

	// FIXME: Move all this together at the top using a defer as done elsewhere.
	//  Maybe we need to declare `res` in the signature.
	if span.IsRecordingEvents() {
		span.AddAttributes(
			trace.Int64Attribute("resp_status", int64(res.Status)),
			trace.StringAttribute("msg", res.ErrorMessage),
			trace.Int64Attribute("chain_len", int64(len(res.Chain))),
		)
	}

	c.peerTracker.logSuccess(peer, build.Clock.Since(connectionStart), uint64(len(res.Chain)))
	// FIXME: We should really log a success only after we validate the response.
	//  It might be a bit hard to do.
	return &res, nil
}
*/