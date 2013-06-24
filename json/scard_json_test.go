package json

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

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
		t.Fail()
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
