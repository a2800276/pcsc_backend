package json

import (
	"testing"
	"strings"
	"bytes"
)

func TestVersion (t *testing.T) {
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
