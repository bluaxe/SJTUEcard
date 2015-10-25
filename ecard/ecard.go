package ecard

import (
	"fmt"
	"github.com/bluaxe/fetch/common"
	"io/ioutil"
	"net"
	"net/url"
	"regexp"
	"strconv"
)

var u urls = NewUrls()

func Login(user common.User, se *Session) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Login failed! [", r, "]")
			ok = false
		}
	}()

	startSession(se)

	getPhoto(se)

	code := getCode(se)

	form := make(url.Values)
	form.Set("name", user.Sid)
	form.Set("userType", "1")
	form.Set("passwd", user.Passwd)
	form.Set("loginType", "2")
	form.Set("rand", code)
	//form.Set("imageField.x", "0")
	//form.Set("imageField.y", "0")
	//fmt.Println(val.Encode())

	resp, err := se.client.PostForm(u.login, form)
	if err != nil || resp.StatusCode != 200 {
		panic(err)
	}

	// fmt.Println(resp.Status)

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	re := regexp.MustCompile("accleft")
	return re.FindStringIndex(string(body)) != nil
}

func initQuery(se *Session) string {
	resp, _ := se.client.Get(u.query)
	body := getBody(resp)
	se.conti = getContinue(body)

	re := regexp.MustCompile("account\" class = .*\n.*\n.*value=\"([^\"]*)\"")
	res := re.FindSubmatch(body)
	account := []byte("")

	if res != nil {
		account = res[1]
		return string(account)
	} else {
		panic("Can not find account !")
	}
}

func Query(se *Session, start, end string) []common.Record {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Query Failed ! [", r, "]")
		}
	}()

	cid := initQuery(se)
	fmt.Printf("Card ID:%q\n", cid)

	//Post Account
	posturl := u.query + "?__continue=" + se.conti

	val := make(url.Values)
	val.Set("account", cid)
	val.Set("inputObject", "all")

	resp, err := se.client.PostForm(posturl, val)
	if err != nil || resp.StatusCode != 200 {
		panic(err)
	}
	body := getBody(resp)
	se.conti = getContinue(body)

	//Post Date
	posturl = u.query + "?__continue=" + se.conti
	form := make(url.Values)
	form.Set("inputStartDate", start)
	form.Set("inputEndDate", end)

	resp, err = se.client.PostForm(posturl, form)
	if err != nil || resp.StatusCode != 200 {
		panic(err)
	}
	body = getBody(resp)
	se.conti = getContinue(body)
	//fmt.Printf("continue=%s", cont)

	//Get Logs
	posturl = u.query + "?__continue=" + se.conti

	form = make(url.Values)
	resp, err = se.client.PostForm(posturl, form)
	if err != nil || resp.StatusCode != 200 {
		panic(err)
	}
	body = getBody(resp)

	num := pageNum(body)
	//fmt.Printf("总共: %d页\n", num)

	recs := records(body)

	pform := make(url.Values)
	for i := 2; i <= num; i++ {
		pform.Set("pageNum", strconv.Itoa(i))
		resp, _ = se.client.PostForm(u.pages, pform)
		if err != nil || resp.StatusCode != 200 {
			panic(err)
		}
		body = getBody(resp)
		recs = append(recs, records(body)...)
	}
	// fmt.Println(recs)
	return recs
}

func startSession(se *Session) {
	se.client.Get(u.home)
}

func getPhoto(se *Session) {
	se.client.Get(u.pwdphoto)
}

func getCode(se *Session) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Get Code failed !", r)
		}
	}()
	resp, _ := se.client.Get(u.chkpic)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	soc := "127.0.0.1:30196"
	tcpAddr, _ := net.ResolveTCPAddr("tcp", soc)
	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()

	length := strconv.Itoa(len(body))
	//fmt.Printf(length)

	_, _ = conn.Write([]byte(length))
	_, _ = conn.Write(body)

	reply := make([]byte, 32)
	_, _ = conn.Read(reply)
	res := string(reply)

	fmt.Printf("code is :[" + res + "]\n")

	//writeFile("1.jpg", body)

	return res[0:4]
}
