package json

import (
	"encoding/json"
	"bytes"
	"io"
	"fmt"
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

func decodeFully(r io.Reader, into interface{})(err error) {
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

func encodeError(mes string, w io.Writer)(err error) {
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

	buffer2:= bytes.NewBuffer(buffer.Bytes())


	var message = ScardRequest{}
	if err = decodeFully(&buffer, &message); err != nil {
		return
	}
	switch (message.Method) {
		case "version":
			return ScardVersion(buffer2, w)
		default:
		println(err.Error())
			return encodeError(fmt.Sprintf("unknown method: %s", message.Method),w) 
	}
}

func ScardVersion(r io.Reader, w io.Writer) (err error) {
	req := ScardRequest{}

	if err = decodeFully(r, &req); err != nil {
		return
	}

	switch req.Method {
		case "version":
			res := ScardVersionResponse{}
			res.Error="0"
			res.Version=scard.Version()
			encoder := json.NewEncoder(w)
			return encoder.Encode(res)
		default:
			return encodeError(fmt.Sprintf("incorrect method: %s", req.Method), w)
	}

}

// 
// 
// // establishContext
// {
// 	method: "establishContext"
// }
// {
// 	error: 0
// 	ctx: 1234
// }	
// // releaseContext
// {
// 	method: "establishContext"
// 	ctx: 
// }
// {
// 	error: 0
// }	
// // cancel
// // valid
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
