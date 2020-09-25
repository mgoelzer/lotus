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
	mandatoryPaymentIntervalInBytes         = 1048576  // 1 mb
	mandatoryPaymentIntervalIncreaseInBytes = 10485760 // 10 mb
)

// Protocol ID used on the wire
const (
	SimpleRetieveProtocolID = "/fil/simple-retrieve/0.0.1"
)

// Request/Response Types
const (
	// ReqRespFoo = N    // ==> int N corresponds to RequestFoo, RespnoseFoo structs
	ReqRespInitialize            = 1
	ReqRespConfirmTransferParams = 2
	ReqRespTransfer              = 3
	ReqRespVoucher               = 4
)

// -- Response Codes --
const (
	// Commonly used response codes
	ResponseCodeOk             = 0
	ResponseCodeGeneralFailure = 1

	// Add'l error states for Initialize
	ResponseCodeInitializeNoCid = 101

	// Add'l error states for ConfirmTransferParams
	ResponseCodeConfirmTransferParamsWrongParams = 201

	// Add'l error states for Transfer
	// (none)

	// Add'l error states for Voucher
	ResponseCodeVoucherSigInvalid = 301
)

// TODO: make this enum a constant in this code
///
// SignedVoucher:
//  enum SignatureType {
//    Secp256k1 = 1,
//    BLS = 2,
//}
///

//
// Deserialization structs
//
type GenericRequestOrResponse struct {
	ReqOrResp    string `json:"type"`
	Request      int
	Response     int
	ResponseCode int
	ErrorMessage string
}

type GenericRequest struct {
	ReqOrResp string `json:"type"`
	Request   int
}

//
// "Initialize" roundtrip structs
//
type RequestInitialize struct {
	GenericRequest
	PchAddr string
	Cid     string // TODO:  use native go-cid Cid type
	Offset0 int64
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
