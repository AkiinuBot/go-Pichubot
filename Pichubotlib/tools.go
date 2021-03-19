package tools

import (
	"fmt"
	"regexp"
)

// in this package you can use tools to help your code

func CQcodeParse(rawmessage string) []map[string]interface{} {
	var re = regexp.MustCompile(`(?m)\[CQ\:(.+?)\]`)
	var output []map[string]interface{}
	for i, match := range re.FindAllString(rawmessage, -1) {
		match
		fmt.Println(match, "found at index", i)

	}
	return output
}
