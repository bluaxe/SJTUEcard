package main

import (
	"fmt"
	"github.com/bluaxe/SJTUEcard/common"
	"github.com/bluaxe/SJTUEcard/ecard"
	"github.com/bluaxe/SJTUEcard/persist"
	"github.com/bluaxe/SJTUEcard/server"
	"github.com/bluaxe/SJTUEcard/service"
	"time"
)

func ecardtest() {
	user := common.User{
		Sid:    "115033910041",
		Passwd: "171289",
	}
	se := ecard.NewSession()

	ok := ecard.Login(user, &se)
	if ok {
		psn := "axe:axe@tcp(neo.bile.dog:3306)/axe?charset=utf8"
		persist.Init(psn)
		db := persist.GetDB(psn)
		defer persist.ReleaseDB(db)
		// persist.CreateRecordTable(db)

		recs := ecard.Query(&se, "20150301", "20150331")
		persist.InsertRecords(db, recs)

		recs = ecard.Query(&se, "20150401", "20150430")
		persist.InsertRecords(db, recs)

		recs = ecard.Query(&se, "20150501", "20150531")
		persist.InsertRecords(db, recs)

		recs = ecard.Query(&se, "20150601", "20150630")
		persist.InsertRecords(db, recs)

		recs = ecard.Query(&se, "20150701", "20150730")
		persist.InsertRecords(db, recs)

		recs = ecard.Query(&se, "20150801", "20150830")
		persist.InsertRecords(db, recs)

		recs = ecard.Query(&se, "20150901", "20150931")
		persist.InsertRecords(db, recs)

		recs = ecard.Query(&se, "20151001", "20151031")
		persist.InsertRecords(db, recs)

	}
	// persist.DumpPlaceMap()
}

func dbtest() {
	persist.Test("axe:axe@tcp(neo.bile.dog:3306)/axe")
}

func servertest() {
	server.Start("0.0.0.0:8000")
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	// dbtest()
	// ecardtest()
	// servertest()
	psn := "root:blue@tcp(axe.so:3306)/axe?charset=utf8"
	persist.Init(psn)

	service.StartFetcher(time.Second * 30)
	server.Start("0.0.0.0:8000")
}
