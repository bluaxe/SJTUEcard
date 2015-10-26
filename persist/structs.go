package persist

type sqls struct {
	CREATE_RECORD_TABLE string
	INSERT_TO_RECORD    string

	CREATE_PLACE_TABLE string
	FIND_PLACE_ID      string
	INSERT_TO_PLACE    string
	SELECT_ALL_PLACES  string

	SELECT_RECENT_TOP string

	CREATE_USER_TABLE  string
	INSERT_TO_USER     string
	SET_USER_VALID     string
	GET_USER           string
	GET_ALL_VALID_USER string
	UPDATE_USER_PWD    string
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

	s.SELECT_RECENT_TOP = `select places.name, am from (
		    	select place, sum(ammount) as am from records where sid=? and date >= '20151001' group by place
		) stat
		left join places on places.id= stat.place order by am;`

	s.CREATE_USER_TABLE = `CREATE TABLE users(
			sid VARCHAR(16),
			passwd VARCHAR(32) NOT NULL,
			name VARCHAR(64) NOT NULL DEFAULT '',
			valid DECIMAL(1) NOT NULL DEFAULT 0,
			PRIMARY KEY(sid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8`

	s.INSERT_TO_USER = `INSERT INTO users(sid, passwd, name) VALUES(?, ?, ?)`
	s.SET_USER_VALID = `UPDATE users SET valid=? where sid=?`
	s.GET_USER = `SELECT sid, passwd, name FROM users where sid=?`
	s.GET_ALL_VALID_USER = `SELECT sid, passwd, name FROM users where valid=1`
	s.UPDATE_USER_PWD = `UPDATE users SET passwd=? where sid=?`

	return s
}
