package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	_ "github.com/go-sql-driver/mysql_original"
)

var db sql.DB
var timeNow string

func main() {

	timeNow = time.Now().Format("2006-01-02")
	fmt.Println(timeNow)

	db, err := sql.Open("mysql", "root:liuliu@tcp(127.0.0.1:3306)/chexiang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for i := 0; i < 7379; i++ {
		//		DelectTables(db,i)
		CreatTables(db, i)

		url := fmt.Sprint("http://car.chexiang.com/product/", i, ".htm")
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("get url err:", err)
			return
		}
		defer resp.Body.Close()

		bs_body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("read resp body err:", err)
			return
		}

		name, price, err := ParseWebBody(bs_body)
		if err != nil {
			fmt.Println("ParseWebBody err:", err)
			return
		}

		remark := ""
		updateMysql(db, i, price, remark)
		fmt.Println(i, name, price, remark)
	}
}

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

func updateMysql(db *sql.DB, id int, price string, remark string) {

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

func ParseWebBody(body []byte) (string, string, error) {
	reg, err := regexp.Compile("car-details.*\n.*<h1>(.*)</h1>")
	if err != nil {
		fmt.Println("Compile err:", err)
		return "", "", err
	}
	results := reg.FindSubmatch(body)
	if len(results) < 2 {
		return "", "", err
	}
	//fmt.Println(string(results[1]))
	name := string(results[1])

	reg, err = regexp.Compile("<strong>&yen;(.*)</strong>")
	if err != nil {
		fmt.Println("Compile err:", err)
		return "", "", err
	}
	results = reg.FindSubmatch(body)
	if len(results) < 2 {
		return "", "", err
	}
	price := string(results[1])

	return name, price, nil
}
