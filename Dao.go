package DhtCrawler

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Dao struct {
	user     string
	password string
	ip       string
	port     int
	database string
	db       *sql.DB
	HashIns1 *sql.Stmt
	HashIns2 *sql.Stmt
}

func NewDao(user string, password string, ip string, port int, database string) *Dao {
	dao := new(Dao)
	dao.user = user
	dao.password = password
	dao.ip = ip
	dao.port = port
	dao.database = database

	dao.Init()
	return dao
}

func (dao *Dao) Init() {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=90s&collation=utf8_general_ci", dao.user, dao.password, dao.ip, dao.port, dao.database)
	db, err := sql.Open("mysql", dns)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	dao.db = db

	HashIns1, err := dao.db.Prepare("INSERT INTO `hash1` (`id`, `hash`) VALUES (NULL, ?);")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	dao.HashIns1 = HashIns1

	HashIns2, err := dao.db.Prepare("INSERT INTO `hash2` (`id`, `hash`) VALUES (NULL, ?);")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	dao.HashIns2 = HashIns2

}
