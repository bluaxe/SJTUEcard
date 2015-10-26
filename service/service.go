package service

import (
	// "database/sql"
	"fmt"
	"github.com/bluaxe/SJTUEcard/common"
	"github.com/bluaxe/SJTUEcard/ecard"
	"github.com/bluaxe/SJTUEcard/persist"
	"time"
)

type serviceProvider func()

var services map[string]serviceProvider

func init() {

}

func doFetchAllRecord(user common.User, se *ecard.Session) {

	db := persist.GetDefaultDB()
	defer persist.ReleaseDB(db)

	recs := ecard.Query(se, "20151001", "20151031")
	persist.InsertRecords(db, recs)

	recs = ecard.Query(se, "20150901", "20150931")
	persist.InsertRecords(db, recs)

	recs = ecard.Query(se, "20150801", "20150830")
	persist.InsertRecords(db, recs)

	recs = ecard.Query(se, "20150701", "20150730")
	persist.InsertRecords(db, recs)

	recs = ecard.Query(se, "20150601", "20150630")
	persist.InsertRecords(db, recs)

	recs = ecard.Query(se, "20150501", "20150531")
	persist.InsertRecords(db, recs)

	recs = ecard.Query(se, "20150401", "20150430")
	persist.InsertRecords(db, recs)

	recs = ecard.Query(se, "20150301", "20150331")
	persist.InsertRecords(db, recs)

	fmt.Println("Fetch Done.")
}

func FetchAllRecord(user common.User) (string, error) {

	se := ecard.NewSession()
	ok := ecard.Login(user, &se)
	if ok {
		go doFetchAllRecord(user, &se)
		return "processing.", nil
	} else {
		return "account not valid.", nil
	}
}

func addAccount(user common.User) {
	persist.InsertUser(user)
	persist.SetUserValid(user, 1)
}

func updateAccount(user common.User) {
	udb, err := persist.GetUser(user.Sid)
	if err != nil {
		fmt.Println("Not fount account ... @ ", err)
		addAccount(user)
		FetchAllRecord(user)
	} else {
		if udb.Passwd != user.Passwd {
			persist.UpdatePwd(user)
			fmt.Println("update user pwd")
		}
		persist.SetUserValid(user, 1)
	}
}

func CheckAccount(user common.User) (string, error) {

	se := ecard.NewSession()

	ok := ecard.Login(user, &se)

	if ok {
		go updateAccount(user)
		return "yes", nil
	} else {
		return "no", nil
	}
}

func StartFetcher(t time.Duration) {
	ticker := time.NewTicker(t)
	go func() {
		for _ = range ticker.C {
			fetchAllUser()
		}
	}()
}

func fetchAllUser() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic in fetchAllUser :", r)
		}
	}()
	fmt.Println("fetch service start")

	users := persist.GetAllValidUsers()
	for _, user := range users {
		fetchBetween(user, "20151024", "20151026")
	}
	fmt.Println("fetch service done")
}

func fetchBetween(user common.User, start, end string) {
	db := persist.GetDefaultDB()
	defer persist.ReleaseDB(db)

	fmt.Println("fetch service [Login] for", user.Sid)
	se := ecard.NewSession()
	ok := ecard.Login(user, &se)
	if ok {
		fmt.Println("fetch service [Fetch] for", user.Sid)
		recs := ecard.Query(&se, start, end)
		persist.InsertRecords(db, recs)
	} else {
		fmt.Println("fetch service [Login Failed] for", user.Sid)
	}
}
