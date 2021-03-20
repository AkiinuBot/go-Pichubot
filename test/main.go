package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	var re = regexp.MustCompile(`(?m)\[CQ\:(.+?)\]`)
	output := make([]map[string]interface{}, 0)
	str := ""
	for _, match := range re.FindAllStringSubmatch(str, -1) {
		split := strings.Split(match[1], ",")
		parsed := make(map[string]interface{})
		for t, stats := range split {
			if t == 0 {
				parsed["name"] = stats
			} else {
				split2 := strings.SplitN(stats, "=", 2)
				parsed[split2[0]] = split2[1]
			}
		}
		output = append(output, parsed)

	}
	fmt.Println(len(output))
}
