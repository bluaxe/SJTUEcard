package ecard

import (
	"net/http"
	"net/http/cookiejar"
)

type Session struct {
	conti  string
	client *http.Client
}

type urls struct {
	login    string
	home     string
	pwdphoto string
	chkpic   string
	query    string
	pages    string
}

func NewUrls() urls {
	var u urls
	u.login = "http://ecard.sjtu.edu.cn/loginstudent.action"
	u.home = "http://ecard.sjtu.edu.cn/homeLogin.action"
	u.pwdphoto = "http://ecard.sjtu.edu.cn/getpasswdPhoto.action"
	u.chkpic = "http://ecard.sjtu.edu.cn/getCheckpic.action?rand=5475.172977894545"
	u.query = "http://ecard.sjtu.edu.cn/accounthisTrjn.action"
	u.pages = "http://ecard.sjtu.edu.cn/accountconsubBrows.action"
	return u
}

func NewSession() Session {
	var s Session

	s.conti = ""

	cookie, err := cookiejar.New(nil)
	if err != nil {
		panic("New cookiejar error! ")
	}

	// agent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.71 Safari/537.36"
	// req.Header.Set("User-Agent", agent)

	s.client = &http.Client{}
	s.client.Jar = cookie

	return s
}
