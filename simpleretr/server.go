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
