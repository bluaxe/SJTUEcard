package main

import (
    "fmt"
	"net/http"
	"net/url"
	"io/ioutil"
	"net/http/cookiejar"
	"net"
	"strconv"
	"os"
	"strings"
	"regexp"
)
func main() {
	do()
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
		getUserInfo(cookie)
	}else{
		fmt.Println("wrong account!")
	}
}

func login(ck *cookiejar.Jar, code string) bool {
	posturl := "http://ecard.sjtu.edu.cn/loginstudent.action"

	val :=make(url.Values)
	val.Set("name", "5110309666")
	val.Set("userType", "1")
	val.Set("passwd", "171289")
	val.Set("loginType", "2")
	val.Set("rand", code[0:4])
	val.Set("imageField.x", "0")
	val.Set("imageField.y", "0")
	//fmt.Println(val.Encode())


	client := &http.Client{
		Jar:ck,
	}

	agent:= "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.71 Safari/537.36"
	bodyType:="application/x-www-form-urlencoded"

	req, _ := http.NewRequest("POST", posturl, strings.NewReader(val.Encode()))
	req.Header.Set("Content-Type", bodyType)
	req.Header.Set("User-Agent", agent)

	resp, err:= client.Do(req)

	//resp, err := client.PostForm(posturl, val)

	if err!=nil {
		panic("Some error occurs")
	}

	//fmt.Println(resp.Status)

	body, err :=ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	re := regexp.MustCompile("accleft")
	return re.FindStringIndex(string(body)) != nil

	//writeFile("out.html", body)
}

func writeFile(filename string, body[]byte){
	f , _ := os.Create(filename)
	defer f.Close()
	f.Write([]byte(body))
}

func getpic(ck *cookiejar.Jar) string {
	geturl := "http://ecard.sjtu.edu.cn/getCheckpic.action?rand=5475.172977894545"
	
	client := &http.Client{
		Jar:ck,
	}

	resp, _:= client.Get(geturl)
	body, _:=ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	soc := "127.0.0.1:30196"
	tcpAddr, _ := net.ResolveTCPAddr("tcp", soc)
	conn, _ := net.DialTCP("tcp",nil, tcpAddr)
	defer conn.Close()

	length := strconv.Itoa(len(body))
	//fmt.Printf(length)

	_,_ = conn.Write([]byte(length))
	_,_ = conn.Write(body)

	reply := make([]byte, 32)
	_, _ = conn.Read(reply)
	res := string(reply)

	fmt.Printf("code is :[" + res+"]\n")

	writeFile("1.jpg", body)	

	return res
}

func getPhoto(ck *cookiejar.Jar) {
	geturl := "http://ecard.sjtu.edu.cn/getpasswdPhoto.action"	
	client := &http.Client{
		Jar:ck,
	}
	resp, _ := client.Get(geturl)
	defer resp.Body.Close()
}

func getUserInfo(ck * cookiejar.Jar) {
	geturl := "http://ecard.sjtu.edu.cn/accountcardUser.action"
	client := &http.Client{
		Jar:ck,
	}
	resp, _ := client.Get(geturl)
	defer resp.Body.Close()

	body,_ := ioutil.ReadAll(resp.Body)
	writeFile("user.html", body)

}

func getSession() *cookiejar.Jar {
	geturl := "http://ecard.sjtu.edu.cn/homeLogin.action"		
	
	cookie, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:cookie,
	}

	resp, _:= client.Get(geturl)
	defer resp.Body.Close()
	
	return cookie
}
