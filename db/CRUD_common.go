package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Pgdb struct {
	db *sql.DB
}

type DbConnPara struct {
	Host     string
	Port     int32
	User     string
	Password string
	Dbname   string
}

func NewPgDb(para *DbConnPara) *Pgdb {

	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		para.Host, para.Port, para.User, para.Password, para.Dbname)

	db, err := sql.Open("postgres", dbInfo)
	if nil != err {
		fmt.Println("Open fail, ", err)
		return nil
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return &Pgdb{db: db}
}

func (p *Pgdb) Close() {
	p.db.Close()
}
