package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"webCrawler/crawler"

	_ "github.com/go-sql-driver/mysql"
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
		crawler.CreatTables(db, i)

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

		name, price, err := crawler.ParseWebBody(bs_body)
		if err != nil {
			fmt.Println("ParseWebBody err:", err)
			return
		}

		remark := ""
		crawler.UpdateMysql(db, i, timeNow, price, remark)
		fmt.Println(i, name, price, remark)
	}
}
