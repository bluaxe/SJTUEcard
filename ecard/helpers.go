package ecard

import (
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/bluaxe/SJTUEcard/common"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func getBody(resp *http.Response) []byte {
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return body
}

func getContinue(body []byte) string {
	re := regexp.MustCompile("\\?__continue=([^\"]*)")
	res := re.FindSubmatch(body)

	cont := []byte("")
	if res != nil {
		cont = res[1]
		//fmt.Printf("%q\n", cont)
		return string(cont)
	} else {
		panic("Can not Find continue!")
	}
}

func records(body []byte) []common.Record {
	records := make([]common.Record, 0)

	re := regexp.MustCompile("<tr class=\"listbg2?\">(.*\n.*){13}.*<")
	res := re.FindAll(body, -1)
	if res != nil {
		// fmt.Printf("length: %d\n", len(res))
		for _, match := range res {
			records = append(records, line(match))
		}
	} else {
		fmt.Println("Not found!")
	}
	return records
}

func gbk2Uni(str []byte) string {
	dec := mahonia.NewDecoder("gbk")
	return strings.Trim(dec.ConvertString(string(str)), " ")
}

func line(line []byte) common.Record {
	var rec common.Record

	re := regexp.MustCompile(".*>(.*)<.*\n.*>(.*)<.*\n.*>(.*)<.*\n.*>(.*)<.*\n.*>(.*)<.*\n.*\n.*>(.*)<.*\n.*>(.*)<.*\n.*>(.*)<.*\n.*>(.*)<.*\n.*>(.*)<")
	res := re.FindAllSubmatch(line, -1)
	if res != nil {
		list := res[0]
		datetime := strings.Split(gbk2Uni(list[1]), " ")
		rec.Date = datetime[0]
		rec.Time = datetime[1]
		//TODO spilt Date to Date and Time
		rec.Sid = gbk2Uni(list[2])
		rec.Username = gbk2Uni(list[3])
		rec.Class = gbk2Uni(list[4])
		rec.Place = gbk2Uni(list[5])
		rec.Ammount, _ = strconv.ParseFloat(gbk2Uni(list[6]), 64)
		rec.Rest, _ = strconv.ParseFloat(gbk2Uni(list[7]), 64)
		rec.Balance, _ = strconv.ParseFloat(gbk2Uni(list[8]), 64)
		rec.Status = gbk2Uni(list[10])

	} else {
		panic(fmt.Sprintf("\tGet Line Content of %s failed !", string(line)))
	}
	return rec
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
		panic("Not Found Page Number!")
	}
}

func writeFile(filename string, body []byte) {
	f, _ := os.Create(filename)
	defer f.Close()
	f.Write([]byte(body))
}
