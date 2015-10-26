package server

import (
	"fmt"
	"github.com/bluaxe/fetch/common"
	// "github.com/bluaxe/fetch/ecard"
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
	body, _ := service.CheckAccount(user)

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

	res, _ := service.FetchAllRecord(user)
	fmt.Fprintf(w, res)
}
