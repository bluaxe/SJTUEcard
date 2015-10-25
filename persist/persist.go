package persist

import (
	"database/sql"
	"fmt"
	"github.com/bluaxe/fetch/common"
	_ "github.com/go-sql-driver/mysql"
)

var s sqls = InitSQLS()
var placemap map[string]common.Place = make(map[string]common.Place)
var pdb *sql.DB
var default_dsn string

func Test(dsn string) {
	fmt.Println(dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("select id,value from axe")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var id int
	var str string
	for rows.Next() {
		err := rows.Scan(&id, &str)
		if err != nil {
			panic(err)
		}
		fmt.Println(id, str)
	}

	// CreateRecordTable(db)
	// CreatePlaceTable(db)

}

func GetDb() *sql.DB {
	return GetDB(default_dsn)
}

func GetDB(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	return db
}

func ReleaseDB(db *sql.DB) {
	db.Close()
}

func InsertRecords(db *sql.DB, recs []common.Record) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
		}
	}()
	stmt, err := db.Prepare(s.INSERT_TO_RECORD)
	if err != nil {
		panic(err)
	}

	for _, rec := range recs {
		p := getPlaceFromName(rec.Place, db)
		_, err := stmt.Exec(rec.Date, rec.Time, rec.Sid, rec.Class, rec.Username, p.Id, rec.Ammount, rec.Balance, rec.Rest, rec.Status)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(rec.Date, rec.Time, rec.Username)
		}
	}

}

func Init(dsn string) {
	default_dsn = dsn
	db := GetDB(dsn)
	defer ReleaseDB(db)
	loadPlaceMap(db)
}

func loadPlaceMap(db *sql.DB) {
	rows, err := db.Query(s.SELECT_ALL_PLACES)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var p common.Place
		err := rows.Scan(&(p.Id), &(p.Name), &(p.Nick))
		if err != nil {
			panic(err)
		}
		placemap[p.Name] = p
	}
}

func DumpPlaceMap() {
	fmt.Println("Dumping --------------------")
	for name, p := range placemap {
		fmt.Println("\t", name, p.Id, p.Nick)
	}
}

func getPlaceFromName(name string, db *sql.DB) common.Place {
	if _, ok := placemap[name]; ok {
		return placemap[name]
	}
	var p common.Place
	p.Name = name
	p.Nick = ""
	insertPlace(db, &p)
	placemap[name] = p
	return p
}

func insertPlace(db *sql.DB, place *common.Place) {
	res, err := db.Exec(s.INSERT_TO_PLACE, place.Name, place.Nick)
	if err != nil {
		panic(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}
	place.Id = id
}

func CreateRecordTable(db *sql.DB) {
	_, err := db.Exec(s.CREATE_RECORD_TABLE)
	if err != nil {
		fmt.Println(err)
	}
}

func CreatePlaceTable(db *sql.DB) {
	_, err := db.Exec(s.CREATE_PLACE_TABLE)
	if err != nil {
		fmt.Println(err)
	}
}
