// Package simpleretr contains the simple-retrieve server.
//
// ** Summary Protocol Description **
//
// simple-retrieve is a libp2p protocol by which a peer (the client) can request a
// Payload CID and receive it in chunks (from the server).  After each chunk is sent
// to the client by the server, the server will pause and wait for client to send a
// payment voucher.  The server will check the voucher, and if it is valid send the
// next chunk of bytes.  This process continues until all bytes are sent, followed
// by the last voucher, at which point the stream can be closed by either side.
//
// ** Detailed Protocol Description **
//
// Here are a few terms we will use to describe the protocol:
//  - Client is the peer requesting data (and providing vouchers), a.k.a. "requester."
//  Note that this package does not implement a Client, only a Server.
//  - Server is the peer providing data (and accepting vouchers), a.k.a. "responder."
//
// Every step of the protocol is a request-response round trip.  This protocol,
// like HTTP, is fundamentally text-based.  Requests and responses are typed structs in Go
// named as RequestFoo and ResponseFoo, but on the wire they are marshalled to JSON objects
// encoded as (potentially very long) UTF8 strings.
//
// Requests look like this on the wire:
//
//  {"type":"requestFoo",
//   "request":"integer-type-of-request-response",
//   "parameter1":"value1", ... }
//
// The integer request/response types are defined in simpleretr/protocol.go and correspond to
// a typed RequestFoo struct in go.  The number parameters and their names are specific
// to that particular RequestFoo.
//
// Responses look like this on the wire:
//
//  {"type":"response",
//  "response":"integer-type-of-request-response",
//  "responseCode":"integer-response-code",
//  "parameter1":"value1", ...      <--- NOT CURRENTLY SUPPORTED
//  "data":"base64-encoded-bytes" }
//
// A response is integer response code, optional parameters, and an optional data block.
// In Go, they are marshalled/unmarshalled to/from structs, which are named as ResponseFoo,
// where Foo is determined by the integer-type-of-request-response field.
//
// The response code is an integer that indicates success, a specific error
// condition, or another status; these codes are defined in this package (protocol.go).
// The data block, when present, contains a base64 encoded string of bytes.  A data block
// can be of any length, but typical sizes would be on the order of single digit MBs up to
// hundreds of megabytes.
//
// (JSON serialization was chosen in order to (1) make this protocol easily
// interoperable with clients written in other languages, particularly Javascript,
// and (2) enable easier inspection on the wire.)
//
// Below is the basic, no-error flow of the protocol.  Refer to protocol.go to see the
// definitions of the Go structs corresponding to these JSON request/response strings.
//
//  - Client connects to Server and opens a stream for this protocol.
//  - Client sends struct RequestInitialize and server responds with ResponseInitialize
//  -
//
// TODO: finish this ^^
//
package simpleretr
