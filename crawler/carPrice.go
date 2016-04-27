package crawler

import (
	"database/sql"
	"fmt"
)

func DelectTables(db *sql.DB, id int) {

	_, err := db.Exec(fmt.Sprint("drop table if exists id_", id, ";"))
	if err != nil {
		fmt.Println("updateMysql delect", err)
	}
}

func CreatTables(db *sql.DB, id int) {

	_, err := db.Exec(fmt.Sprint("create table if not exists id_", id, "(time date primary KEY, price char(15),remark varchar(100));"))
	if err != nil {
		fmt.Println("updateMysql Create", err)
	}
}

func UpdateMysql(db *sql.DB, id int, timeNow string, price string, remark string) {

	stmt, err := db.Prepare(fmt.Sprint("INSERT id_", id, " SET time=?,price=?,remark=?"))
	defer stmt.Close()
	if err != nil {
		fmt.Println("updateMysql prepare", err)
	}

	_, err = stmt.Exec(timeNow, price, remark)
	if err != nil {
		fmt.Println("updateMysql execute", err)
	}
}
