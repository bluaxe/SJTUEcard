package service

import (
	"fmt"
	"github.com/bluaxe/fetch/common"
	"github.com/bluaxe/fetch/ecard"
	"github.com/bluaxe/fetch/persist"
)

type serviceProvider func()

var services map[string]serviceProvider

func init() {

}

func FetchAllRecord(user common.User, se *ecard.Session) {
	db := persist.GetDefaultDB()
	defer persist.ReleaseDB(db)
	// persist.CreateRecordTable(db)

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
