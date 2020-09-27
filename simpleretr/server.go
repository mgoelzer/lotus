package simpleretr

import (
	//"bufio"
	//"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	//"go.opencensus.io/trace"
	//"golang.org/x/xerrors"

	//"github.com/ipfs/go-cid"
	"bufio"
	"bytes"
	b64 "encoding/base64"
	"io"

	//cborutil "github.com/filecoin-project/go-cbor-util"
	"github.com/libp2p/go-libp2p-core/network"
)

// server implements exchange.Server. It services requests for the
// libp2p ChainExchange protocol.
type server struct {
}

var _ Server = (*server)(nil)

// Global variable for now -- MockDataStore is just temp
// (Not protected by lock because MockDataStore is read-only and thread safe)
var ds MockDataStore = NewMockDataStore()

// Per connection state
// No lock because the protocol is inherently serial and single threaded
type connectionState struct {
	Cid     string
	Offset0 int64
}

func NewServer() Server {
	return &server{}
}

func (s *server) HandleStream(stream network.Stream) {
	defer stream.Close()
	fmt.Println(">>> [start] HandleStream()")

	// Instantiate the state struct for this specific /fil/simple-retrieve connection
	cstate := &connectionState{}

	for {
		err, requestJson := getIncomingJsonString(bufio.NewReader(stream))
		if err != nil {
			log.Warnf("failed to read incoming message: %s\n", err)
			break
			//return
		}
		fmt.Printf("READ> %s\n\n", requestJson)
		UnmarshallJsonAndHandle(requestJson, stream, cstate)
	}

	fmt.Println(">>> [end] HandleStream()")

	//if err := cborutil.WriteCborRPC(stream, resp); err != nil {}

	/*	var req Request
		if err := cborutil.ReadCborRPC(bufio.NewReader(stream), &req); err != nil {
			log.Warnf("failed to read block sync request: %s", err)
			return
		}
		log.Infow("block sync request",
			"start", req.Head, "len", req.Length)

		resp, err := s.processRequest(ctx, &req)
		if err != nil {
			log.Warn("failed to process request: ", err)
			return
		}

		_ = stream.SetDeadline(time.Now().Add(WriteResDeadline))
		if err := cborutil.WriteCborRPC(stream, resp); err != nil {
			_ = stream.SetDeadline(time.Time{})
			log.Warnw("failed to write back response for handle stream",
				"err", err, "peer", stream.Conn().RemotePeer())
			return
		}
		_ = stream.SetDeadline(time.Time{})
	*/
}

func UnmarshallJsonAndHandle(jsonStr string, stream network.Stream, cstate *connectionState) error {
	genericReqOrResp := GenericRequestOrResponse{}
	if err := json.Unmarshal([]byte(jsonStr), &genericReqOrResp); err != nil {
		return err
	}
	if genericReqOrResp.ReqOrResp == "request" {
		switch genericReqOrResp.Request {
		case ReqRespInitialize:
			reqInitialize := RequestInitialize{}
			if err := json.Unmarshal([]byte(jsonStr), &reqInitialize); err != nil {
				return err
			}
			if err := HandleRequestInitialize(&reqInitialize, stream, cstate); err != nil {
				return err
			}
		case ReqRespConfirmTransferParams:
			fmt.Println("[sretrieve] ConfirmTransferParams")
		case ReqRespTransfer:
			fmt.Println("[sretrieve] Transfer")
			reqTransfer := RequestTransfer{}
			log.Infof("[sretrieve] Received ")
			if err := json.Unmarshal([]byte(jsonStr), &reqTransfer); err != nil {
				return err
			}
			if err := HandleRequestTransfer(&reqTransfer, stream, cstate); err != nil {
				return err
			}
		case ReqRespVoucher:
			fmt.Println("[sretrieve] Voucher")
		}
	} else {
		return errors.New("[sretrieve] Ignoring: server should never receive a Response struct")
	}
	return nil
}

func HandleRequestInitialize(reqInitialize *RequestInitialize, stream network.Stream, cstate *connectionState) error {
	fmt.Println("[sretrieve] --RequestInitialize--")
	fmt.Printf("[sretrieve] .ReqOrResp = %v\n", reqInitialize.ReqOrResp)
	fmt.Printf("[sretrieve] .Request   = %v\n", reqInitialize.Request)
	fmt.Printf("[sretrieve] .PchAddr = %v\n", reqInitialize.PchAddr)
	fmt.Printf("[sretrieve] .Cid = %v\n", reqInitialize.Cid)
	fmt.Printf("[sretrieve] .Offset0 = %v\n", reqInitialize.Offset0)

	// Assert Offset0==0 in this version

	// Make sure we have this Cid + Cid's total size at Offset0
	cid := reqInitialize.Cid
	hasCid, err := ds.HasCid(cid)
	if err != nil || !hasCid {
		return errors.New(fmt.Sprintf("Server does not have cid '%v'", cid))
	}

	// Cid - save it somewhere
	cstate.Cid = cid
	cstate.Offset0 = reqInitialize.Offset0

	// PchAddress - save it somewhere for later voucher validation

	// Prepare a ResponseInitialize struct

	// Serialize to Json and fire away
	jsonStr := fmt.Sprintf(`{"type":"response","response":%v,"responseCode":%v,"totalBytes":68157440}`, ReqRespInitialize, ResponseCodeOk)
	err = writeToStream(stream, jsonStr)
	return err
}

func HandleRequestTransfer(reqTransfer *RequestTransfer, stream network.Stream, cstate *connectionState) error {
	fmt.Println("[sretrieve] --RequestTransfer--")
	fmt.Printf("[sretrieve] .ReqOrResp = %v\n", reqTransfer.ReqOrResp)
	fmt.Printf("[sretrieve] .Request   = %v\n", reqTransfer.Request)
	fmt.Printf("[sretrieve] .N         = \"%v\"\n", reqTransfer.N)
	fmt.Printf("[sretrieve] .Offset    = %v\n", reqTransfer.Offset)

	// Get N and Offset
	Nstr := reqTransfer.N
	var N int64
	N, err := strconv.ParseInt(Nstr, 10, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("[sretrieve] (HandleRequestTransfer) N from request could not be converted to int64 (N=\"%v\")", reqTransfer.N))
	}
	var Offset int64
	Offset = int64(reqTransfer.Offset)
	fmt.Printf("[sretrieve] (HandleRequestTransfer) N=%d (int64), Offset=%d (int64)\n", N, Offset)

	// Get N bytes starting at Offset
	var bytes []byte
	bytes, err = ds.GetBytes(cstate.Cid, Offset, N)
	if err != nil {
		return errors.New(fmt.Sprintf("[sretrieve] (HandleRequestTransfer) GetBytes failed with '%v'", err))
	}
	var bytesBase64 string
	bytesBase64 = b64.StdEncoding.EncodeToString(bytes)

	// Serialize to Json and fire away
	jsonStr := fmt.Sprintf(`{"type":"response","response":%v,"responseCode":%v,"data":"%v"}`, ReqRespTransfer, ResponseCodeOk, bytesBase64)
	fmt.Printf("[sretrieve] (HandleRequestTransfer) json response: \"%v\"\n", jsonStr)
	err = writeToStream(stream, jsonStr)
	return err
}

func getIncomingJsonString(r io.Reader) (error, string) {
	var err error
	intermediateBuffer := bytes.Buffer{}
	buf := make([]byte, 32)
	for {
		n, err := r.Read(buf)
		//fmt.Printf("n = %v | buf = %v | buf[:n]= %q\n", n, buf, buf[:n])
		intermediateBuffer.Write(buf[:n])
		if err == io.EOF {
			fmt.Printf("readIncomingString:  EOF\n")
			break
		} else if err != nil {
			fmt.Printf("readIncomingString:  error while reading (%v)\n", err)
			break
		}
	}
	return err, intermediateBuffer.String()
}

func writeToStream(stream network.Stream, s string) error {
	w := bufio.NewWriter(stream)
	sBuf := []byte(s)
	sBytes := len(s)
	fmt.Printf("writeToStream: sending back '%s'\n", s)

	//n, err := w.WriteString(s)
	n, err := w.Write(sBuf)
	if err != nil {
		return err
	}

	if err = w.Flush(); err != nil {
		return err
	}

	if n == sBytes {
		log.Infof("writeToStream:  all bytes written to stream (%v bytes)", n)
		return nil
	} else {
		log.Errorf("Wrong number of bytes sent back to client (%v bytes send, %v bytes expectd)\n", sBytes, n)
		return errors.New("Error(writeToStream):  not all bytes written to stream")
	}
}

// TODO:  Delete this, just around for what I can steal from blocksync
/*
// Validate and service the request. We return either a protocol
// response or an internal error.
func (s *server) processRequest(ctx context.Context, req *Request) (*Response, error) {
	validReq, errResponse := validateRequest(ctx, req)
	if errResponse != nil {
		// The request did not pass validation, return the response
		//  indicating it.
		return errResponse, nil
	}

	return s.serviceRequest(ctx, validReq)
}

// Validate request. We either return a `validatedRequest`, or an error
// `Response` indicating why we can't process it. We do not return any
// internal errors here, we just signal protocol ones.
func validateRequest(ctx context.Context, req *Request) (*validatedRequest, *Response) {
	_, span := trace.StartSpan(ctx, "chainxchg.ValidateRequest")
	defer span.End()

	validReq := validatedRequest{}

	validReq.options = parseOptions(req.Options)
	if validReq.options.noOptionsSet() {
		return nil, &Response{
			Status:       BadRequest,
			ErrorMessage: "no options set",
		}
	}

	validReq.length = req.Length
	if validReq.length > MaxRequestLength {
		return nil, &Response{
			Status: BadRequest,
			ErrorMessage: fmt.Sprintf("request length over maximum allowed (%d)",
				MaxRequestLength),
		}
	}
	if validReq.length == 0 {
		return nil, &Response{
			Status:       BadRequest,
			ErrorMessage: "invalid request length of zero",
		}
	}

	if len(req.Head) == 0 {
		return nil, &Response{
			Status:       BadRequest,
			ErrorMessage: "no cids in request",
		}
	}
	validReq.head = types.NewTipSetKey(req.Head...)

	// FIXME: Add as a defer at the start.
	span.AddAttributes(
		trace.BoolAttribute("blocks", validReq.options.IncludeHeaders),
		trace.BoolAttribute("messages", validReq.options.IncludeMessages),
		trace.Int64Attribute("reqlen", int64(validReq.length)),
	)

	return &validReq, nil
}

func (s *server) serviceRequest(ctx context.Context, req *validatedRequest) (*Response, error) {
	_, span := trace.StartSpan(ctx, "chainxchg.ServiceRequest")
	defer span.End()

	chain, err := collectChainSegment(s.cs, req)
	if err != nil {
		log.Warn("block sync request: collectChainSegment failed: ", err)
		return &Response{
			Status:       InternalError,
			ErrorMessage: err.Error(),
		}, nil
	}

	status := Ok
	if len(chain) < int(req.length) {
		status = Partial
	}

	return &Response{
		Chain:  chain,
		Status: status,
	}, nil
}

func collectChainSegment(cs *store.ChainStore, req *validatedRequest) ([]*BSTipSet, error) {
	var bstips []*BSTipSet

	cur := req.head
	for {
		var bst BSTipSet
		ts, err := cs.LoadTipSet(cur)
		if err != nil {
			return nil, xerrors.Errorf("failed loading tipset %s: %w", cur, err)
		}

		if req.options.IncludeHeaders {
			bst.Blocks = ts.Blocks()
		}

		if req.options.IncludeMessages {
			bmsgs, bmincl, smsgs, smincl, err := gatherMessages(cs, ts)
			if err != nil {
				return nil, xerrors.Errorf("gather messages failed: %w", err)
			}

			// FIXME: Pass the response to `gatherMessages()` and set all this there.
			bst.Messages = &CompactedMessages{}
			bst.Messages.Bls = bmsgs
			bst.Messages.BlsIncludes = bmincl
			bst.Messages.Secpk = smsgs
			bst.Messages.SecpkIncludes = smincl
		}

		bstips = append(bstips, &bst)

		// If we collected the length requested or if we reached the
		// start (genesis), then stop.
		if uint64(len(bstips)) >= req.length || ts.Height() == 0 {
			return bstips, nil
		}

		cur = ts.Parents()
	}
}

func gatherMessages(cs *store.ChainStore, ts *types.TipSet) ([]*types.Message, [][]uint64, []*types.SignedMessage, [][]uint64, error) {
	blsmsgmap := make(map[cid.Cid]uint64)
	secpkmsgmap := make(map[cid.Cid]uint64)
	var secpkincl, blsincl [][]uint64

	var blscids, secpkcids []cid.Cid
	for _, block := range ts.Blocks() {
		bc, sc, err := cs.ReadMsgMetaCids(block.Messages)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		// FIXME: DRY. Use `chain.Message` interface.
		bmi := make([]uint64, 0, len(bc))
		for _, m := range bc {
			i, ok := blsmsgmap[m]
			if !ok {
				i = uint64(len(blscids))
				blscids = append(blscids, m)
				blsmsgmap[m] = i
			}

			bmi = append(bmi, i)
		}
		blsincl = append(blsincl, bmi)

		smi := make([]uint64, 0, len(sc))
		for _, m := range sc {
			i, ok := secpkmsgmap[m]
			if !ok {
				i = uint64(len(secpkcids))
				secpkcids = append(secpkcids, m)
				secpkmsgmap[m] = i
			}

			smi = append(smi, i)
		}
		secpkincl = append(secpkincl, smi)
	}

	blsmsgs, err := cs.LoadMessagesFromCids(blscids)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	secpkmsgs, err := cs.LoadSignedMessagesFromCids(secpkcids)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return blsmsgs, blsincl, secpkmsgs, secpkincl, nil
}
*/
