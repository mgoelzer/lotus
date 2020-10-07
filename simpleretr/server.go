package simpleretr

import (
	//"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"time"

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

	// For waiting on go routine completion
	//var wg sync.WaitGroup
	//wg.Add(1)

	// Instantiate the state struct for this specific /fil/simple-retrieve connection
	cstate := &connectionState{}

	//go func(cstate *connectionState, stream network.Stream, wg *sync.WaitGroup) {
	//	defer wg.Done()
	for {
		fmt.Printf("[sretrieve] (HandleStream) (%s) Inside for loop\n", time.Now().String())
		fmt.Printf("[sretrieve] (HandleStream) (%s) Starting wait on `bufio.NewReader(*stream)`\n", time.Now().String())
		rdr := bufio.NewReader(stream)
		fmt.Printf("[sretrieve] (HandleStream) (%s) Finished wait on `bufio.NewReader(*stream)`\n", time.Now().String())
		err, requestJson := getIncomingJsonString(rdr)
		fmt.Printf("[sretrieve] (HandleStream) requestJson = %s, err = %v\n", requestJson, err)
		if err == io.EOF {
			fmt.Printf("[sretrieve] (HandleStream) Note:  getIncomingJsonString returned EOF; continuing...\n\n")
		}
		bNullTerm := requestJson[len(requestJson)-1] == '\x00'
		fmt.Printf("[sretrieve] (HandleStream) len(requestJson) = %v, bNullTerm=%v \n\n", len(requestJson), bNullTerm)
		if bNullTerm {
			requestJson = requestJson[:len(requestJson)-1]
		}
		fmt.Printf("[sretrieve] (HandleStream) len(requestJson) = %v, bNullTerm=%v \n\n", len(requestJson), bNullTerm)

		if len(requestJson) > 0 {
			fmt.Printf("[sretrieve] (HandleStream) <<< READ >>> %s\n\n", requestJson)
			UnmarshallJsonAndHandle(requestJson, stream, cstate)
		} else if err != nil {
			log.Warnf("[sretrieve] (HandleStream) failed to read incoming message: %s\n", err)
			break
			//return
		}
		time.Sleep(1 * time.Second)
	}
	//}(cstate, stream, &wg)

	fmt.Printf("[sretrieve] (HandleStream) Waiting on wait group\n")
	//wg.Wait()
	fmt.Printf("[sretrieve] (HandleStream) Finished waiting on wait group\n")
	fmt.Println("[sretrieve] end - HandleStream()")
}

func UnmarshallJsonAndHandle(jsonStr string, stream network.Stream, cstate *connectionState) error {
	genericReqOrResp := GenericRequestOrResponse{}
	fmt.Printf("[sretrieve] UnmarshallJsonAndHandle:  entered function")
	if err := json.Unmarshal([]byte(jsonStr), &genericReqOrResp); err != nil {
		log.Errorf("[sretrieve] UnmarshallJsonAndHandle:  exiting: failed to unmarshal jsonStr='%s'", jsonStr)
		return err
	}
	if genericReqOrResp.ReqOrResp == "request" {
		switch genericReqOrResp.Request {
		case ReqRespInitialize:
			fmt.Printf("[sretrieve] UnmarshallJsonAndHandle:  in case ReqRespInitialize")
			reqInitialize := RequestInitialize{}
			if err := json.Unmarshal([]byte(jsonStr), &reqInitialize); err != nil {
				log.Errorf("[sretrieve] UnmarshallJsonAndHandle:  json.Unmarshal failed (err='%v')", err)
				return err
			}
			if err := HandleRequestInitialize(&reqInitialize, stream, cstate); err != nil {
				log.Errorf("[sretrieve] UnmarshallJsonAndHandle:  HandleRequestInitialize failed (err='%v')", err)
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
		default:
			log.Errorf("[sretrieve] UnmarshallJsonAndHandle:  fell through to default case (jsonStr='%s')", jsonStr)
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
	// TODO

	// Prepare a ResponseInitialize struct
	// TODO

	// Serialize to Json and fire away
	// TODO:  fake response
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
		fmt.Printf("[sretrieve] getIncomingJsonString:  n = %v \n                                    buf = %v \n                                    buf[:n]= %s\n", n, buf, buf[:n])
		intermediateBuffer.Write(buf[:n])
		if err == io.EOF {
			fmt.Printf("[sretrieve] getIncomingJsonString:  EOF\n")
			break
		} else if buf[n-1] == '\x00' {
			fmt.Printf("[sretrieve] getIncomingJsonString:  NULL terminator\n")
			break
		} else if err != nil {
			fmt.Printf("[sretrieve] getIncomingJsonString:  error while reading (%v)\n", err)
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
