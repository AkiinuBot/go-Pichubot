package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	var re = regexp.MustCompile(`(?m)\[CQ\:(.+?)\]`)
	var output []map[string]interface{}
	str := "[CQ:at,qq=123456=2]"
	for _, match := range re.FindAllStringSubmatch(str, -1) {
		split := strings.Split(match[1], ",")
		for t, v := range split {
			fmt.Println(v)
			if t == 1 {
				_type =
			} else {

			}
		}

	}
}
