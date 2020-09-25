package simpleretr

import (
	//"bufio"
	//"context"
	"encoding/json"
	"errors"
	"fmt"

	//"go.opencensus.io/trace"
	//"golang.org/x/xerrors"

	//"github.com/ipfs/go-cid"
	"bufio"
	"bytes"
	"io"

	//cborutil "github.com/filecoin-project/go-cbor-util"
	inet "github.com/libp2p/go-libp2p-core/network"
)

// server implements exchange.Server. It services requests for the
// libp2p ChainExchange protocol.
type server struct {
}

var _ Server = (*server)(nil)

func NewServer() Server {
	return &server{}
}

func (s *server) HandleStream(stream inet.Stream) {
	defer stream.Close()
	fmt.Println(">>> [start] HandleStream()")

	//for
	{
		// err, requestJson := getIncomingJsonString(bufio.NewReader(stream))
		// if err != nil {
		// 	log.Warnf("failed to read incoming message: %s\n", err)
		// 	//break
		// 	return
		// }

		// //else if err == io.EOF { // Need to wait on stream instead of busy wait with sleep
		// //	time.Sleep(5 * time.Second)
		// //}

		//fmt.Printf("READ> %s\n\n", requestJson)
		//UnmarshallJsonAndHandle(requestJson, stream)

		w := bufio.NewWriter(stream)
		str := fmt.Sprintf(`{"type":"response","response":%v,"responseCode”:%v,"totalBytes":68157440}`, ReqRespInitialize, ResponseCodeOk)
		strBytes := len(str)

		//////
		// _, err = w.WriteString(str)
		// if err != nil {
		// 	log.Warnf("failed to write to stream with WriteString: %s\n", err)
		// 	return
		// }
		//////
		buf := []byte(str)
		n, err := w.Write(buf)
		if err != nil {
			log.Warnf("failed to write to stream with Write: %s\n", err)
			return
		}
		if n == strBytes {
			log.Infof("Correct number of bytes sent back to client (%v bytes)\n", n)
		} else {
			log.Errorf("Wrong number of bytes sent back to client (%v bytes send, %v bytes expectd)\n", n, strBytes)
		}
		//////
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

func UnmarshallJsonAndHandle(jsonStr string, stream inet.Stream) error {
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
			if err := HandleRequestInitialize(&reqInitialize, stream); err != nil {
				return err
			}
		case ReqRespConfirmTransferParams:
			fmt.Println("ConfirmTransferParams")
		case ReqRespTransfer:
			fmt.Println("Transfer")
		case ReqRespVoucher:
			fmt.Println("Voucher")
		}
	} else {
		return errors.New("Server should never receive a Response struct")
	}
	return nil
}

func HandleRequestInitialize(reqInitialize *RequestInitialize, stream inet.Stream) error {
	fmt.Println("--RequestInitialize--")
	fmt.Printf(".ReqOrResp = %v\n", reqInitialize.ReqOrResp)
	fmt.Printf(".Request   = %v\n", reqInitialize.Request)
	fmt.Printf(".PchAddr = %v\n", reqInitialize.PchAddr)
	fmt.Printf(".Cid = %v\n", reqInitialize.Cid)
	fmt.Printf(".Offset0 = %v\n", reqInitialize.Offset0)

	// Assert Offset0==0 in this version

	// Make sure we have this Cid + Cid's total size at Offset0

	// PchAddress - save it somewhere for later voucher validation

	// Prepare a ResponseInitialize struct

	// Serialize to Json and fire away

	w := bufio.NewWriter(stream)
	str := fmt.Sprintf(`{"type":"response","response":%v,"responseCode”:%v,"totalBytes":68157440}`, ReqRespInitialize, ResponseCodeOk)
	strBytes := len(str)
	fmt.Printf("HandleRequestInitialize: sending back '%s'\n", str)
	n, err := w.WriteString(str)
	if err != nil {
		return err
	}
	if n == strBytes {
		fmt.Println("HandleRequestInitialize:  all bytes written to stream")
		return nil
	} else {
		return errors.New("Error(HandleRequestInitialize):  not all bytes written to stream")
	}
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
