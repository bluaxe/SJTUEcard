package server

import (
	"fmt"
	"net/http"
	"net/url"
)

var entrys map[string]handler = make(map[string]handler)

func init() {
	entrys["test"] = testfunc
	entrys["chk"] = checkAccount
}

func testfunc(v url.Values, w http.ResponseWriter) {
	fmt.Fprintf(w, "test func\n")
}

func dispatcher(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form.Encode())
	en := r.Form.Get("action")
	f, have := entrys[en]
	if have && en != "" {
		f(r.Form, w)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}
