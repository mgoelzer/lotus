package simpleretr

import (
	//"time"

	//"github.com/filecoin-project/lotus/build"
	//"github.com/filecoin-project/lotus/chain/store"

	//"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log"
	//"golang.org/x/xerrors"

	//"github.com/filecoin-project/lotus/chain/types"
)

var log = logging.Logger("simple-retrieve")

// If you are implementing the other side of this protocol in another language, 
// you should translate this const block into your language.
const (
	// Protocol ID used on the wire
	SimpleRetieveProtocolID = "/fil/simple-retrieve/0.0.1"

	// Request types
	RequestTypeInitialize = 1
)

// TODO: Rename. Make private.
type RequestInitialize struct {
	// List of ordered CIDs comprising a `TipSetKey` from where to start
	// fetching backwards.
	// FIXME: Consider using `TipSetKey` now (introduced after the creation
	//  of this protocol) instead of converting back and forth.
	//Head []cid.Cid
	// Number of block sets to fetch from `Head` (inclusive, should always
	// be in the range `[1, MaxRequestLength]`).
	//Length uint64
	// Request options, see `Options` type for more details. Compressed
	// in a single `uint64` to save space.
	//Options uint64
}

// TODO: Rename. Make private.
type ResponseInitialize struct {
	ResponseCode status
	ErrorMessage string
	Data string
}

type status uint64

const (
	Ok status = 0

	// TODO:  customize errors
	// Errors
	NotFound      = 201
	GoAway        = 202
	InternalError = 203
	BadRequest    = 204
)

/* // Convert status to internal error.
func (res *Response) statusToError() error {
	switch res.Status {
	case Ok, Partial:
		return nil
		// FIXME: Consider if we want to not process `Partial` responses
		//  and return an error instead.
	case NotFound:
		return xerrors.Errorf("not found")
	case GoAway:
		return xerrors.Errorf("not handling 'go away' chainxchg responses yet")
	case InternalError:
		return xerrors.Errorf("block sync peer errored: %s", res.ErrorMessage)
	case BadRequest:
		return xerrors.Errorf("block sync request invalid: %s", res.ErrorMessage)
	default:
		return xerrors.Errorf("unrecognized response code: %d", res.Status)
	}
} */
