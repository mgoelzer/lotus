package simpleretr

import (
	//"time"

	//"github.com/filecoin-project/lotus/build"
	//"github.com/filecoin-project/lotus/chain/store"

	//"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log"
	//"golang.org/x/xerrors"
)

var log = logging.Logger("simple-retrieve")

////////////////////////////////////////////////////////
//                                                    //
//                    CONSTANTS                       //
//                                                    //
////////////////////////////////////////////////////////

// If you are implementing the other side of this protocol in another language,
// translate all the const blocks below to your language.

// Mandatory values (temporary - in future versions nodes can negotiate these)
const (
	mandatoryPaymentIntervalInBytes         = 1024 * 1024     // 1MB
	mandatoryPaymentIntervalIncreaseInBytes = 5 * 1024 * 1024 // 5MB
)

// Protocol ID used on the wire
const (
	SimpleRetieveProtocolID = "/fil/simple-retrieve/0.0.1"
)

// Request/Response Types
const (
	// RequestTypeNumFoo = N     // N corresponds to RequestFoo, RespnoseFoo structs
	ReqRespNumInitialize = 1 // corresponds to RequestInitialize, ResponseInitialize
)

// -- Response Codes --
const (
	// Commonly used response codes
	ResponseCodeOk = /*responseCode*/ 0

	// Error state response codes (Initialize)
	ResponseCodeInitializeNoCid = 101
	ResponseCodeInitializeFail  = 102
)

//
// Initialize roundtrip
//
type RequestInitialize struct {
	CID                            string `json:"cid"` // TODO:  use native go-cid Cid type
	PaymentIntervalInBytes         int64  `json:"paymentIntervalInBytes"`
	PaymentIntervalIncreaseInBytes int64  `json:"paymentIntervalIncreaseInBytes"`
}
type ResponseInitialize struct {
	ResponseCode int
	ErrorMessage string
	Data         string
}

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
