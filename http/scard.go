package http

import (
	"net/http"
	"strings"
	"io"
)

import "github.com/ebfe/go.pcsclite/scard"

type ScardCookie struct {
	cookie *http.Cookie
	ctx    *scard.Context
}

type ScardHandler struct {
	cookies map[string]*ScardCookie
}

func (hdlr *ScardHandler) establishContext(w http.ResponseWriter, req *http.Request) {
	// check cookie set
	// set cookie
	// establish context (have some sort of timeout function synced to cookie timeout)
	// put cookie in map
}

//isvalid
//listReaders
//cancel

func (hdlr *ScardHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	api_func := path[strings.LastIndex(path, "/") : ]
	switch api_func {
		case "EstablishContext":
			hdlr.establishContext(w,req)
		default:
			w.Write([]byte(api_func))
			io.Copy(w, req.Body)
	}
}

