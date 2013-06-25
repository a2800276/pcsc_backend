package json

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

import "github.com/ebfe/go.pcsclite/scard"

func TestVersion(t *testing.T) {
	req := `{
	"method": "version"
	}`
	reader := strings.NewReader(req)
	writer := &bytes.Buffer{}

	if err := ScardJson(reader, writer); err != nil {
		t.Fail()
	} else {
		resp := ScardVersionResponse{}
		if err = decodeFully(writer, &resp); err != nil {
			t.Fail()
		} else {
			if resp.Version != "1.7.4" {
				println(resp.Version)
				t.Fail()
			}
		}
	}
}

func TestCtx(t *testing.T) {
	if ctx, err := getContext(); err != nil {
		t.Fail()
	} else {
		t.Log(ctx)
		_ = contexts[ctx].Release()
	}
}

func getContext() (c Context, err error) {

	req := `{
	"method": "establishContext"
	}`
	reader := strings.NewReader(req)
	writer := &bytes.Buffer{}

	if err = ScardJson(reader, writer); err != nil {
		return
	} else {
		resp := ScardContextResponse{}
		if err = decodeFully(writer, &resp); err != nil {
			return
		}
		return resp.Ctx, nil
	}
}
func TestRelease(t *testing.T) {
	var ctx Context
	var err error

	if ctx, err = getContext(); err != nil {
		t.Errorf("couldn't get initial ctx: %s", err.Error())
	}

	req := `{
	"method": "releaseContext",
	"ctx": "%s"
	}`
	req = fmt.Sprintf(req, ctx)
	reader := strings.NewReader(req)
	writer := &bytes.Buffer{}

	if err := ScardJson(reader, writer); err != nil {
		t.Fail()
	} else {
		resp := ScardResponse{}
		if err = decodeFully(writer, &resp); err != nil {
			t.Fail()
		} else {
			if resp.Error != "0" {
				t.Logf("unexpected error: %s", resp.Error)
				t.FailNow()
			}
			if contexts[ctx] != nil {
				t.Errorf("context still stored after release")
			}
		}
	}
}

func TestValid(t *testing.T) {
	req := `{
	"method": "isValid",
	"ctx": "%s"
	}`

	var ctx Context
	var err error

	if ctx, err = getContext(); err != nil {
		t.Errorf("couldn't get initial ctx: %s", err.Error())
	}

	req = fmt.Sprintf(req, ctx)
	reader := strings.NewReader(req)
	writer := &bytes.Buffer{}

	if err := ScardJson(reader, writer); err != nil {
		t.Fail()
	} else {
		resp := ScardResponse{}
		if err = decodeFully(writer, &resp); err != nil {
			t.Fail()
		} else {
			if resp.Error != "0" {
				t.Logf("unexpected error: %s", resp.Error)
				t.FailNow()
			}
		}
	}
	contexts[ctx].Release()
	// check again with released context

	reader = strings.NewReader(req)
	writer = &bytes.Buffer{}
	if err := ScardJson(reader, writer); err != nil {
		t.Fail()
	} else {
		resp := ScardResponse{}
		if err = decodeFully(writer, &resp); err != nil {
			t.Fail()
		} else {
			if resp.Error != "INVALID_HANDLE" {
				t.Logf("unexpected error: %s", resp.Error)
				t.FailNow()
			}
		}
	}

}

func TestListReaders(t *testing.T) {
	req := `{
	"method": "listReaders",
	"ctx": "%s"
	}`

	var ctx Context
	var err error

	if ctx, err = getContext(); err != nil {
		t.Errorf("couldn't get initial ctx: %s", err.Error())
	}

	req = fmt.Sprintf(req, ctx)
	reader := strings.NewReader(req)
	writer := &bytes.Buffer{}

	if err := ScardJson(reader, writer); err != nil {
		t.Fatal(err)
	} else {
		resp := ScardListReadersResponse{}
		if err = decodeFully(writer, &resp); err != nil {
			t.Fail()
		} else {
			if resp.Error != "0" {
				t.Logf("unexpected error: %s", resp.Error)
				t.FailNow()
			}
			t.Logf("found %d readers\n", len(resp.Readers))
			for _, rdr := range resp.Readers {
				t.Logf("reader: %s\n", rdr)
			}
		}
	}
	contexts[ctx].Release()
	// check again with released context

	reader = strings.NewReader(req)
	writer = &bytes.Buffer{}
	if err := ScardJson(reader, writer); err != nil {
		t.Fail()
	} else {
		resp := ScardResponse{}
		if err = decodeFully(writer, &resp); err != nil {
			t.Fail()
		} else {
			if resp.Error != "INVALID_HANDLE" {
				t.Logf("unexpected error: %s", resp.Error)
				t.FailNow()
			}
		}
	}

}

func TestConnect(t *testing.T) {
	var reader string
	if ctx, err := scard.EstablishContext(); err != nil {
		t.Fatal(err)
	} else {
		defer ctx.Release()
		if rdrs, err := ctx.ListReaders(); err != nil {
			t.Fatal(err)
		} else {
			contexts["123"] = ctx
			reader = rdrs[0]
		}
	}

	req := `{
		"method":"connect",
		"ctx": "123",
		"reader": "%s",
		"shareMode": "EXCLUSIVE",
		"protocol": "ANY"
	}`

	req = fmt.Sprintf(req, reader)
	rder := strings.NewReader(req)
	writer := &bytes.Buffer{}

	if err := ScardJson(rder, writer); err != nil {
		t.Fatal(err)
	} else {
		resp := ScardConnectResponse{}
		if err = decodeFully(writer, &resp); err != nil {
			t.Fatal(err)
		} else {
			if resp.Error != "0" {
				t.Logf("unexpected error: %s", resp.Error)
				t.FailNow()
			}
			t.Logf("card: %s\n", resp.Card)
			if cards[resp.Card] == nil {
				t.Error("card not in server")
			}
			if err = cards[resp.Card].Disconnect(scard.UNPOWER_CARD); err != nil {
				t.Fatal(err)
			} else {
			}
			cards[resp.Card] = nil
		}
	}

}

func TestStatus(t *testing.T) {
	var reader string
	if ctx, err := scard.EstablishContext(); err != nil {
		t.Fatal(err)
	} else {
		defer ctx.Release()
		if rdrs, err := ctx.ListReaders(); err != nil {
			t.Fatal(err)
		} else {
			contexts["123"] = ctx
			reader = rdrs[0]
		}
	}

	var card *scard.Card
	var err error
	if card, err = contexts["123"].Connect(reader, scard.SHARE_EXCLUSIVE, scard.PROTOCOL_ANY); err != nil {
		t.Fatal(err)
	}

	cards[Card("123")] = card
	req := `{
			"method":"status",
			"card":"123"
		}`

	req = fmt.Sprintf(req, reader)
	rder := strings.NewReader(req)
	writer := &bytes.Buffer{}

	if err := ScardJson(rder, writer); err != nil {
		t.Fatal(err)
	} else {
		resp := ScardStatusResponse{}
		if err = decodeFully(writer, &resp); err != nil {
			t.Fatal(err)
		} else {
			if resp.Error != "0" {
				t.Logf("unexpected error: %s", resp.Error)
				t.FailNow()
			}
			t.Logf("card: %s\n", resp.Card)
			t.Logf("reader: %s (%s)\n", resp.Reader, reader)
			t.Logf("state: %s\n", resp.State)
			t.Logf("proto: %s\n", resp.ActiveProtocol)
			t.Logf("atr: %s\n", resp.ATR)

			if cards[resp.Card] == nil {
				t.Error("card not in server")
			}
			cards[resp.Card].Disconnect(scard.UNPOWER_CARD)
			cards[resp.Card] = nil
		}
	}

}

func TestDisconnect(t *testing.T) {
	var reader string
	if ctx, err := scard.EstablishContext(); err != nil {
		t.Fatal(err)
	} else {
		defer ctx.Release()
		if rdrs, err := ctx.ListReaders(); err != nil {
			t.Fatal(err)
		} else {
			contexts["123"] = ctx
			reader = rdrs[0]
		}
	}

	var card *scard.Card
	var err error
	if card, err = contexts["123"].Connect(reader, scard.SHARE_EXCLUSIVE, scard.PROTOCOL_ANY); err != nil {
		t.Fatal(err)
	}

	cards[Card("123")] = card

	req := `{
		"method":"disconnect",
		"card":"123",
		"disposition":"UNPOWER_CARD"
	}`

	req = fmt.Sprintf(req, reader)
	rder := strings.NewReader(req)
	writer := &bytes.Buffer{}

	if err := ScardJson(rder, writer); err != nil {
		t.Fatal(err)
	} else {
		resp := ScardResponse{}
		if err = decodeFully(writer, &resp); err != nil {
			t.Fatal(err)
		} else {
			if resp.Error != "0" {
				t.Fatalf("unexpected error: %s", resp.Error)
			}
			if cards["123"] != nil {
				t.Fatal("card still in server")
			}
		}
	}
}

func TestTransmit(t *testing.T) {
	var reader string
	if ctx, err := scard.EstablishContext(); err != nil {
		t.Fatal(err)
	} else {
		defer ctx.Release()
		if rdrs, err := ctx.ListReaders(); err != nil {
			t.Fatal(err)
		} else {
			contexts["123"] = ctx
			reader = rdrs[0]
		}
	}

	var card *scard.Card
	var err error
	if card, err = contexts["123"].Connect(reader, scard.SHARE_EXCLUSIVE, scard.PROTOCOL_ANY); err != nil {
		t.Fatal(err)
	}
	// currently stupid workaround so SCM SCR 3310 will work
	// it seems to me that this reader expects the card to be reinserted...
	// Disconnecting the card with disposition = reset won't work.
	if err = card.Reconnect(scard.SHARE_EXCLUSIVE, scard.PROTOCOL_ANY, scard.RESET_CARD); err != nil {
		t.Fatal(err)
	}
	cards[Card("123")] = card
	defer card.Disconnect(scard.UNPOWER_CARD)

	// use any old credit card to test
	pse := []byte("1PAY.SYS.DDF01\x00")
	select_apdu := append([]byte{0x00, 0xA4, 0x04, 0x00, 0x0e}, pse...)

	apdu := hex.EncodeToString(select_apdu)

	req := `{
		"method":"transmit",
		"card":"123",
		"data":"%s"
	}`

	req = fmt.Sprintf(req, apdu)
	rder := strings.NewReader(req)
	writer := &bytes.Buffer{}
	if err = card.BeginTransaction(); err != nil {
		t.Fatal(err)
	}
	if err := ScardJson(rder, writer); err != nil {
		t.Fatal(err)
	} else {
		resp := ScardTransmitResponse{}
		if err = decodeFully(writer, &resp); err != nil {
			t.Fatal(err)
		} else {
			if resp.Error != "0" {
				t.Fatalf("unexpected error: %s", resp.Error)
			}
			t.Logf("response: %s", resp.Data)
		}
	}
	if err = card.EndTransaction(scard.LEAVE_CARD); err != nil {
		t.Fatal(err)
	}
	contexts["123"] = nil
	cards[Card("123")] = nil
}
