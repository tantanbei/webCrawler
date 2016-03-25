package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"

	_ "github.com/go-sql-driver/mysql_original"
)

type Cars struct {
	carsMap map[string]string
	sync.Mutex
}

func (self *Cars) init() {
	self.carsMap = make(map[string]string)
}

func (self *Cars) UpdateMap(name string, price string) {
	self.Lock()
	self.carsMap[name] = price
	self.Unlock()
}

func main() {
	cars := &Cars{}
	cars.init()
	//UrlArray := make([]string, 0)
	for i := 0; i < 7379; i++ {
		//		DelectTables(i)
		CreatTables(i)

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

		updateMysql(i, price)
		fmt.Println(i, name, price)
	}

}

func DelectTables(id int) {
	db, err := sql.Open("mysql", "root:lbx@tcp(127.0.0.1:3306)/chexiang")
	defer db.Close()
	if err != nil {
		fmt.Println("updateMysql open", err)
	}

	_, err = db.Exec(fmt.Sprint("drop table if exists id_", id, ";"))
	if err != nil {
		fmt.Println("updateMysql delect", err)
	}
}

func CreatTables(id int) {
	db, err := sql.Open("mysql", "root:lbx@tcp(127.0.0.1:3306)/chexiang")
	defer db.Close()
	if err != nil {
		fmt.Println("updateMysql open", err)
	}

	_, err = db.Exec(fmt.Sprint("create table if not exists id_", id, "(time date primary KEY, price char(15),remark varchar(100));"))
	if err != nil {
		fmt.Println("updateMysql Create", err)
	}

}

func updateMysql(id int, price string) {
	db, err := sql.Open("mysql", "root:lbx@tcp(127.0.0.1:3306)/chexiang")
	defer db.Close()
	if err != nil {
		fmt.Println("updateMysql open", err)
	}

	stmt, err := db.Prepare(fmt.Sprint("INSERT id_", id, " SET time=?,price=?,remark=?"))
	defer stmt.Close()
	if err != nil {
		fmt.Println("updateMysql prepare", err)
	}

	_, err = stmt.Exec("2016-03-23", price, "")
	if err != nil {
		fmt.Println("updateMysql execute", err)
	}
}

func (self *Cars) run(url string) {
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
	self.UpdateMap(name, price)
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
