package main

import (
	"database/sql"
	"fmt"
)

//create the all name of cars on chexiang, and support chinese
func CreatNameTable(db *sql.DB) {
	_, err := db.Exec(fmt.Sprint("create table if not exists chexiangCarName(id int primary KEY not null, name varchar(100) not null)ENGINE=InnoDB DEFAULT CHARSET=utf8;"))
	if err != nil {
		fmt.Println("chexiangCarName Create", err)
	}
}

//add the chexiang car name into table
func updateNameTable(db *sql.DB, id int, name string) {
	stmt, err := db.Prepare(fmt.Sprint("INSERT chexiangCarName SET id=?, name=?"))
	defer stmt.Close()
	if err != nil {
		fmt.Println("updateMysql prepare", err)
	}

	_, err = stmt.Exec(id, name)
	if err != nil {
		fmt.Println("updateMysql execute", err)
	}
}
