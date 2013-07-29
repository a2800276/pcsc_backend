package main

import (
	"net/http"
	emvhttp "emv/http" 
	"fmt"
	"html"
)
import "github.com/ebfe/go.pcsclite/scard"
import "github.com/a2800276/gaml"


func main() {
	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q\n", html.EscapeString(r.URL.Path))
		if ctx, err := scard.EstablishContext(); err != nil {
			fmt.Fprintf(w, "Something went wrong: %s\n", err.Error())
		} else {
			defer ctx.Release()
			if rdrs, err2 := ctx.ListReaders(); err2 != nil {
				fmt.Fprintf(w, "Could not get readers: %s\n", err2.Error())
			} else {
				for i, rdr := range rdrs {
					fmt.Fprintf(w, "%d %s\n", i, rdr)
				}
			}
		}
	})
	handler := &emvhttp.ScardHandler{}
	http.Handle("/scard/", handler)

	gamlHandler, err := gaml.NewGamlHandlerWithRenderer("./assets/gaml", gaml.DefaultToStringRenderer)
	if err != nil {
		println(err.Error())
	}
	println("here")
	http.Handle("/", gamlHandler)
	http.Handle("/js/", http.FileServer(http.Dir("./assets")))
	http.Handle("/html/", http.FileServer(http.Dir("./assets")))
	http.ListenAndServe(":8080", nil)
	println("here")
}
