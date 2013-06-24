package json

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
)

import "github.com/ebfe/go.pcsclite/scard"

// need to map:

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
	case "listReaders":
		return ScardListReaders(buffer2, w)
	case "connect":
		return ScardConnect(buffer2, w)
	case "status":
		return ScardStatus(buffer2, w)
	case "disconnect":
		return ScardDisconnect(buffer2, w)
	default:
		return encodeError(fmt.Sprintf("unknown method: %s", message.Method), w)
	}
}

type scardFunc func(r io.Reader, w io.Writer) (err error)

func scardTemplate(method string, f scardFunc, r io.Reader, w io.Writer) (err error) {
	req := ScardRequest{}

	if err = decodeFully(r, &req); err != nil {
		return
	}

	switch req.Method {
	case method:
		return f(r, w)
	default:
		return encodeError(fmt.Sprintf("incorrect method: %s", req.Method), w)
	}
}

func ScardVersion(r io.Reader, w io.Writer) (err error) {
	f := func(r io.Reader, w io.Writer) (err error) {
		res := ScardVersionResponse{}
		res.Error = "0"
		res.Version = scard.Version()
		encoder := json.NewEncoder(w)
		return encoder.Encode(res)
	}
	return scardTemplate("version", f, r, w)
}

var contexts = make(map[Context]*scard.Context)
var cards = make(map[Card]*scard.Card)

// generates string tokens representing cards and contexts
func genToken() string {
	return strconv.FormatInt(time.Now().UnixNano(), 16)
}

func ScardEstablishContext(r io.Reader, w io.Writer) (err error) {

	f := func(r io.Reader, w io.Writer) (err error) {
		if ctx, serr := scard.EstablishContext(); serr != nil {
			return encodeError(serr.Error(), w)
		} else {
			res := ScardContextResponse{}
			res.Error = "0"

			token := Context(genToken())
			contexts[token] = ctx
			res.Ctx = token

			encoder := json.NewEncoder(w)
			return encoder.Encode(res)
		}
	}
	return scardTemplate("establishContext", f, r, w)

}

type scardCtxFunc func(ctx *scard.Context, ctx_token Context, w io.Writer) (err error)

func scardCtxTemplate(method string, f scardCtxFunc, r io.Reader, w io.Writer) (err error) {
	req := ScardCtxRequest{}

	if err = decodeFully(r, &req); err != nil {
		return
	}

	switch req.Method {
	case method:
		scard_ctx := contexts[Context(req.Ctx)]
		if scard_ctx == nil {
			return encodeError("unknown ctx", w)
		}
		return f(scard_ctx, req.Ctx, w)
	default:
		return encodeError(fmt.Sprintf("incorrect method: %s", req.Method), w)
	}
}
func ScardReleaseContext(r io.Reader, w io.Writer) (err error) {
	f := func(ctx *scard.Context, tok Context, w io.Writer) (err error) {
		if err = ctx.Release(); err != nil {
			return encodeError(err.Error(), w)
		} else {
			contexts[tok] = nil
			resp := ScardResponse{"0"}
			encoder := json.NewEncoder(w)
			return encoder.Encode(resp)
		}
	}
	return scardCtxTemplate("releaseContext", f, r, w)
}
func ScardIsValid(r io.Reader, w io.Writer) (err error) {
	f := func(ctx *scard.Context, _ Context, w io.Writer) (err error) {
		if valid, err2 := ctx.IsValid(); err != nil {
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

func checkContext(ctx *scard.Context, w io.Writer) (valid bool, err error) {
	if ctx == nil {
		return false, encodeError("UNKNOWN_CTX", w)
	}
	if valid, err = ctx.IsValid(); err != nil {
		return false, encodeError(err.Error(), w)
	} else if !valid {
		return false, encodeError("INVALID_HANDLE",w)
	}
	return
}

func ScardListReaders(r io.Reader, w io.Writer) (err error) {
	f := func(ctx *scard.Context, _ Context, w io.Writer) (err error) {
		var valid bool
		if valid, err = checkContext(ctx, w); !valid {
			return
		}
		if readers, err2 := ctx.ListReaders(); err != nil {
			return encodeError(err2.Error(), w)
		} else {
			resp := ScardListReadersResponse{}
			resp.Error = "0"
			resp.Readers = readers
			encoder := json.NewEncoder(w)
			return encoder.Encode(resp)
		}
	}

	return scardCtxTemplate("listReaders", f, r, w)
}

func connect(req *ScardConnectRequest, w io.Writer) (err error) {
	ctx := contexts[req.Ctx]
	var valid bool
	if valid, err = checkContext(ctx, w); !valid {
		return // checkContext already sent the error.
	}
	if !req.Protocol.OK() || !req.ShareMode.OK() {
		return encodeError("INCORRECT_PARAM", w)
	}

	var card *scard.Card
	if card, err = ctx.Connect(req.Reader, req.ShareMode.Scard(), req.Protocol.Scard()); err != nil {
		return encodeError(err.Error(), w)
	} else {
		jsoncard := Card(genToken())
		cards[jsoncard] = card
		resp := ScardConnectResponse{}
		resp.Error = "0"
		resp.Card = jsoncard
		encoder := json.NewEncoder(w)
		return encoder.Encode(resp)
	}
}
func ScardConnect(r io.Reader, w io.Writer) (err error) {
	req := ScardConnectRequest{}

	if err = decodeFully(r, &req); err != nil {
		return
	}

	switch req.Method {
	case "connect":
		return connect(&req, w)
	default:
		return encodeError(fmt.Sprintf("incorrect method: %s", req.Method), w)
	}
	return
}

func checkCard(card Card, w io.Writer) (scard *scard.Card, err error) {
	scard_card := cards[card]
	if scard_card == nil {
		return nil, encodeError("UNKNOWN_CARD", w)
	}
	return scard_card, nil
}
func status(req *ScardStatusRequest, w io.Writer) (err error) {
	var card *scard.Card
	if card, err = checkCard(req.Card, w); card == nil {
		return
	}

	var status *scard.CardStatus
	if status, err = card.Status(); err != nil {
		return encodeError(err.Error(), w)
	}
	resp := ScardStatusResponse{}
	resp.Error = "0"
	resp.Card = req.Card
	resp.Reader = status.Reader
	resp.ActiveProtocol = ProtocolFromScard(status.ActiveProtocol)
	resp.ATR = hex.EncodeToString(status.ATR)
	encoder := json.NewEncoder(w)
	return encoder.Encode(resp)
}
func ScardStatus(r io.Reader, w io.Writer) (err error) {
	req := ScardStatusRequest{}

	if err = decodeFully(r, &req); err != nil {
		return
	}

	switch req.Method {
	case "status":
		return status(&req, w)
	default:
		return encodeError(fmt.Sprintf("incorrect method: %s", req.Method), w)
	}
	return
}

func ScardDisconnect(r io.Reader, w io.Writer) (err error) {
	req := ScardDisconnectRequest{}

	if err = decodeFully(r, &req); err != nil {
		return
	}

	switch req.Method {
	case "disconnect":
		if !req.Disposition.OK() {
			return encodeError("INCORRECT_PARAM", w)
		}
		var card *scard.Card 
		if card, err = checkCard(req.Card, w); card == nil {
			return
		}
		if err = card.Disconnect(req.Disposition.Scard()); err!= nil {
			return encodeError(err.Error(), w)
		}
		cards[req.Card] = nil
		resp := ScardResponse{}
		resp.Error="0"
		encoder:= json.NewEncoder(w)
		return encoder.Encode(resp) 
	default:
		return encodeError(fmt.Sprintf("incorrect method: %s", req.Method), w)
	}
}
// // cancel
// // reconnect
// {}
// // transmit
