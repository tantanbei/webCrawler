package crawler

import (
	"fmt"
	"regexp"
)

//return the car name and price
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
