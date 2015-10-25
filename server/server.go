package server

import (
	"fmt"
	"github.com/bluaxe/fetch/common"
	"github.com/bluaxe/fetch/ecard"
	"github.com/bluaxe/fetch/persist"
	"github.com/bluaxe/fetch/service"
	"net/http"
	"net/url"
)

func Start(addr string) {
	http.HandleFunc("/", dispatcher)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func checkAccount(v url.Values, w http.ResponseWriter) {
	sid := v.Get("sid")
	pwd := v.Get("pwd")
	if sid == "" || pwd == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	user := common.User{
		Sid:    sid,
		Passwd: pwd,
	}
	se := ecard.NewSession()

	ok := ecard.Login(user, &se)

	var body string = ""
	if ok {
		body = "yes"
	} else {
		body = "no"
	}
	fmt.Fprintf(w, body)
}

func showRecentTop(v url.Values, w http.ResponseWriter) {
	sid := v.Get("sid")
	if sid == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	s := persist.GetRecentTop(sid)
	fmt.Fprintf(w, s)
}

func fetchAllRecord(v url.Values, w http.ResponseWriter) {
	sid := v.Get("sid")
	pwd := v.Get("pwd")
	if sid == "" || pwd == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	user := common.User{
		Sid:    sid,
		Passwd: pwd,
	}

	se := ecard.NewSession()
	ok := ecard.Login(user, &se)

	if ok {
		go service.FetchAllRecord(user, &se)
		fmt.Fprintf(w, "processing.")
	} else {
		fmt.Fprintf(w, "account not correct.")
	}
}
