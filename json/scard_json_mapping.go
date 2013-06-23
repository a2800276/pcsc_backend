package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
)

import "github.com/ebfe/go.pcsclite/scard"

// need to map:

type ScardRequest struct {
	Method string `json:"method"`
}

type ScardResponse struct {
	Error string `json:"error"`
}

type ScardVersionResponse struct {
	ScardResponse
	Version string `json:"version"`
}

type ScardContextResponse struct {
	ScardResponse
	Ctx string `json:"ctx"`
}

type ScardCtxRequest struct {
	ScardRequest
	Ctx string `json:"ctx"`
}


func decodeFully(r io.Reader, into interface{}) (err error) {
	decoder := json.NewDecoder(r)
	for {
		if err = decoder.Decode(into); err == io.EOF {
			break
		} else {
			if err != nil {
				return nil
			}
		}
	}
	return nil
}

func encodeError(mes string, w io.Writer) (err error) {
	resp := ScardResponse{Error: mes}
	encoder := json.NewEncoder(w)
	if err = encoder.Encode(resp); err != nil {
		return
	}
	return nil
}

func ScardJson(r io.Reader, w io.Writer) (err error) {
	buffer := bytes.Buffer{}
	io.Copy(&buffer, r)

	buffer2 := bytes.NewBuffer(buffer.Bytes())

	var message = ScardRequest{}
	if err = decodeFully(&buffer, &message); err != nil {
		return
	}

	//fmt.Printf(">%v<", message)

	switch message.Method {
	case "version":
		return ScardVersion(buffer2, w)
	case "establishContext":
		return ScardEstablishContext(buffer2, w)
	case "releaseContext":
		return ScardReleaseContext(buffer2, w)
	case "isValid":
		return ScardIsValid(buffer2, w)
	default:
		return encodeError(fmt.Sprintf("unknown method: %s", message.Method), w)
	}
}

type scardFunc func(r io.Reader, w io.Writer)(err error)
func scardTemplate(method string, f scardFunc, r io.Reader, w io.Writer)(err error){
	req := ScardRequest{}

	if err = decodeFully(r, &req); err != nil {
		return
	}

	switch req.Method {
	case method:
		return f(r,w)
	default:
		return encodeError(fmt.Sprintf("incorrect method: %s", req.Method), w)
	}
}

func ScardVersion(r io.Reader, w io.Writer) (err error) {
	f := func(r io.Reader, w io.Writer)(err error){
		res := ScardVersionResponse{}
		res.Error = "0"
		res.Version = scard.Version()
		encoder := json.NewEncoder(w)
		return encoder.Encode(res)
	}
	return scardTemplate("version", f, r, w)
}

var contexts = make(map[string]*scard.Context)

func ScardEstablishContext(r io.Reader, w io.Writer) (err error) {


	f := func(r io.Reader, w io.Writer)(err error){
		if ctx, serr := scard.EstablishContext(); serr != nil {
			return encodeError(serr.Error(), w)
		} else {
			res := ScardContextResponse{}
			res.Error = "0"

			token := strconv.FormatInt(time.Now().UnixNano(), 16)
			contexts[token] = ctx
			res.Ctx = token

			encoder := json.NewEncoder(w)
			return encoder.Encode(res)
		}
	}
	return scardTemplate("establishContext", f, r, w)

}

type scardCtxFunc func(ctx string, r io.Reader, w io.Writer)(err error)
func scardCtxTemplate(method string, f scardCtxFunc, r io.Reader, w io.Writer)(err error){
	req := ScardCtxRequest{}

	if err = decodeFully(r, &req); err != nil {
		return
	}

	switch req.Method {
	case method:
		return f(req.Ctx, r,w)
	default:
		return encodeError(fmt.Sprintf("incorrect method: %s", req.Method), w)
	}
}
func ScardReleaseContext(r io.Reader, w io.Writer) (err error) {
	f := func(ctx string, r io.Reader, w io.Writer)(err error){
		scard_ctx := contexts[ctx]
		if scard_ctx == nil {
			return encodeError("unknown ctx", w)
		}
		if err = scard_ctx.Release(); err != nil {
			return encodeError(err.Error(), w)
		} else {
			contexts[ctx] = nil
			resp := ScardResponse{"0"}
			encoder := json.NewEncoder(w)
			return encoder.Encode(resp)
		}
	}
	return scardCtxTemplate("releaseContext", f, r, w)
}
func ScardIsValid(r io.Reader, w io.Writer) (err error) {
	f := func(ctx string, r io.Reader, w io.Writer)(err error){
		scard_ctx := contexts[ctx]
		if scard_ctx == nil {
			return encodeError("unknown ctx", w)
		}
		if valid, err2 := scard_ctx.IsValid(); err != nil {
			return encodeError(err2.Error(), w)
		} else {
			if !valid {
				return encodeError("INVALID_HANDLE", w)
			} else {
				return encodeError("0", w)
			}
		}
	}
	return scardCtxTemplate("isValid", f, r, w)
}

// // cancel
// // listReaders
// {
//
// }
// // listReaders
// {
// 	method: "listReaders"
// 	ctx:
// }
// {
// 	error:
// 	readers: ["",...]
// }
// // connect
// {
// 	method: connect
// 	ctx:
// 	mode
// 	protocol
// }
// {
// 	error
// 	ctx
// 	card
// }
// // reconnect
// {}
// // disconnect
// // transmit
