package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	// var command = "sender"
	// var m string
	// if time.Now().Format("PM") == "PM" {
	// 	m = "下午"
	// } else {
	// 	m = "上午"
	// }
	// fmt.Println(time.Now().Format(fmt.Sprintf("2006年1月2日%s3:04", m)))
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println(strconv.FormatInt(time.Now().UnixNano(), 10))
		}()
	}
	select {}
}
