package main

import (
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/bluaxe/fetch/common"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	a := common.User{
		Sid:    "5110309666",
		Passwd: "171289",
	}
	fmt.Println(a)
	//do()
}

func do() {
	cookie := getSession()
	//ecardurl,_ := url.Parse("http://ecard.sjtu.edu.cn/homeLogin.action")
	//_ := cookie.Cookies(ecardurl)
	//fmt.Printf(ck[0].String())

	getPhoto(cookie)

	code := getpic(cookie)

	ok := login(cookie, code)

	if ok {
		fmt.Println("ok!")
		//getUserInfo(cookie)
		initQuery(cookie)
	} else {
		fmt.Println("wrong account!")
	}
}

func login(ck *cookiejar.Jar, code string) bool {
	posturl := "http://ecard.sjtu.edu.cn/loginstudent.action"

	val := make(url.Values)
	val.Set("name", "5110309666")
	val.Set("userType", "1")
	val.Set("passwd", "171289")
	val.Set("loginType", "2")
	val.Set("rand", code[0:4])
	val.Set("imageField.x", "0")
	val.Set("imageField.y", "0")
	//fmt.Println(val.Encode())

	client := &http.Client{
		Jar: ck,
	}

	agent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.71 Safari/537.36"
	bodyType := "application/x-www-form-urlencoded"

	req, _ := http.NewRequest("POST", posturl, strings.NewReader(val.Encode()))
	req.Header.Set("Content-Type", bodyType)
	req.Header.Set("User-Agent", agent)

	resp, err := client.Do(req)

	//resp, err := client.PostForm(posturl, val)

	if err != nil {
		panic("Some error occurs")
	}

	//fmt.Println(resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	re := regexp.MustCompile("accleft")
	return re.FindStringIndex(string(body)) != nil

	//writeFile("out.html", body)
}

func writeFile(filename string, body []byte) {
	f, _ := os.Create(filename)
	defer f.Close()
	f.Write([]byte(body))
}

func getpic(ck *cookiejar.Jar) string {
	geturl := "http://ecard.sjtu.edu.cn/getCheckpic.action?rand=5475.172977894545"

	client := &http.Client{
		Jar: ck,
	}

	resp, _ := client.Get(geturl)
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

	return res
}

func getPhoto(ck *cookiejar.Jar) {
	geturl := "http://ecard.sjtu.edu.cn/getpasswdPhoto.action"
	client := &http.Client{
		Jar: ck,
	}
	resp, _ := client.Get(geturl)
	defer resp.Body.Close()
}

func getUserInfo(ck *cookiejar.Jar) {
	geturl := "http://ecard.sjtu.edu.cn/accountcardUser.action"
	client := &http.Client{
		Jar: ck,
	}
	resp, _ := client.Get(geturl)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	re := regexp.MustCompile("名：</div>.*\n.*\n.*\n.*<div align=\"left\">(.*)</div>")

	res := re.FindSubmatch(body)
	name := res[1]
	fmt.Printf("%q\n", name)

	//writeFile("user.html", body)
}

func initQuery(ck *cookiejar.Jar) (string, string) {
	client := &http.Client{
		Jar: ck,
	}

	//Init
	geturl := "http://ecard.sjtu.edu.cn/accounthisTrjn.action"
	resp, _ := client.Get(geturl)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	cont := getContinue(body)

	re := regexp.MustCompile("account\" class = .*\n.*\n.*value=\"([^\"]*)\"")
	res := re.FindSubmatch(body)
	account := []byte("")

	if res != nil {
		account = res[1]
		fmt.Printf("ID:%q\n", account)
	}

	//Post Account
	posturl := "http://ecard.sjtu.edu.cn/accounthisTrjn.action?__continue=" + string(cont)
	val := make(url.Values)
	val.Set("account", string(account))
	val.Set("inputObject", "all")

	resp, _ = client.PostForm(posturl, val)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	cont = getContinue(body)

	//Post Date
	posturl = "http://ecard.sjtu.edu.cn/accounthisTrjn.action?__continue=" + string(cont)
	form := make(url.Values)
	form.Set("inputStartDate", "20150601")
	form.Set("inputEndDate", "20150630")

	resp, _ = client.PostForm(posturl, form)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	cont = getContinue(body)
	//fmt.Printf("continue=%s", cont)

	//Get Logs
	posturl = "http://ecard.sjtu.edu.cn/accounthisTrjn.action?__continue=" + string(cont)
	form = make(url.Values)

	resp, _ = client.PostForm(posturl, form)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	num := pageNum(body)
	fmt.Printf("总共: %d页\n", num)

	records(body)

	pform := make(url.Values)
	posturl = "http://ecard.sjtu.edu.cn/accountconsubBrows.action"
	for i := 2; i <= num; i++ {
		pform.Set("pageNum", strconv.Itoa(i))
		resp, _ = client.PostForm(posturl, pform)
		body, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		records(body)
	}

	//writeFile("query.html", body)

	return string(cont), string(account)

}

func pageNum(body []byte) int {
	dec := mahonia.NewDecoder("gbk")
	bodys := dec.ConvertString(string(body))
	re := regexp.MustCompile("共(\\d)页")
	res := re.FindSubmatch([]byte(bodys))
	if res != nil {
		//fmt.Printf("%q\n", res[0])
		num, _ := strconv.Atoi(string(res[1]))
		return num
	} else {
		return 0
	}
}

func records(body []byte) {
	re := regexp.MustCompile("<tr class=\"listbg2?\">(.*\n.*){13}.*<")

	res := re.FindAll(body, -1)
	if res != nil {
		fmt.Printf("length: %d\n", len(res))
		for _, match := range res {
			line(match)
		}
	} else {
		fmt.Println("Not found!")
	}
}

func line(line []byte) {
	re := regexp.MustCompile(".*>(.*)<.*\n.*>(.*)<.*\n.*>(.*)<.*\n.*>(.*)<.*\n.*>(.*)<.*\n.*\n.*>(.*)<.*\n.*>(.*)<.*\n.*>(.*)<.*\n.*>(.*)<.*\n.*>(.*)<")
	res := re.FindAllSubmatch(line, -1)
	dec := mahonia.NewDecoder("gbk")
	if res != nil {
		for _, list := range res {
			for i := 1; i < len(list); i++ {
				match := strings.Trim(string(list[i]), " ")
				fmt.Printf("\t%q", dec.ConvertString(match))
			}
			fmt.Println("")
		}
	}
}

func getContinue(body []byte) string {
	re := regexp.MustCompile("\\?__continue=([^\"]*)")
	res := re.FindSubmatch(body)

	cont := []byte("")
	if res != nil {
		cont = res[1]
		//fmt.Printf("%q\n", cont)
	}
	return string(cont)
}

func getSession() *cookiejar.Jar {
	geturl := "http://ecard.sjtu.edu.cn/homeLogin.action"

	cookie, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookie,
	}

	resp, _ := client.Get(geturl)
	defer resp.Body.Close()

	return cookie
}
