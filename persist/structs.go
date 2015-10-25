package persist

type sqls struct {
	CREATE_RECORD_TABLE string
	INSERT_TO_RECORD    string
	CREATE_PLACE_TABLE  string
	FIND_PLACE_ID       string
	INSERT_TO_PLACE     string
	SELECT_ALL_PLACES   string
}

func InitSQLS() sqls {
	var s sqls
	s.CREATE_RECORD_TABLE = `CREATE TABLE records (
				date    	DATE,
				time    	TIME,
				sid      	VARCHAR(32),
				class   	VARCHAR(32),
				username 	VARCHAR(32),
				place    	int NOT NULL,
				ammount  	DECIMAL(10,2),
				balance  	DECIMAL(10,2),
				rest     	DECIMAL(10,2),
				status   	VARCHAR(32),
				PRIMARY KEY(date, time, sid, ammount)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8`
	s.INSERT_TO_RECORD = `INSERT INTO 
		records(date, time, sid, class, username, place, ammount, balance, rest, status) 
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	s.CREATE_PLACE_TABLE = `CREATE TABLE places(
			id int NOT NULL AUTO_INCREMENT,
			name VARCHAR(64) NOT NULL,
			nick VARCHAR(64) NOT NULL DEFAULT '',
			PRIMARY KEY(id),
			UNIQUE(name)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8`

	s.INSERT_TO_PLACE = `INSERT INTO places(name, id) VALUES(?,?)`

	s.SELECT_ALL_PLACES = `SELECT id,name,nick FROM places`
	return s
}
