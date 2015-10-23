package common

type User struct {
	Username string
	Passwd   string
	Cid      string
	Sid      string
}

type Record struct {
	Date     string
	Time     string
	Username string
	Place    string
	Ammount  float64
	Balance  float64
	Rest     float64
	Status   string
}
