package pichumod

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// in this package you can use tools to help your code

var Connect *websocket.Conn

func init() {
	LinkLog()
}

// 传入带CQ码的消息(string) 解析为[map[]...]并返回;
// 如果消息不带CQ码，则返回空值
func CQcodeParse(rawmessage string) []map[string]interface{} {
	var re = regexp.MustCompile(`(?m)\[CQ\:(.+?)\]`)
	output := make([]map[string]interface{}, 0)
	for _, match := range re.FindAllStringSubmatch(rawmessage, -1) {
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
	switch len(output) {
	case 0:
		return nil
	default:
		return output
	}
}

var LogFile *os.File

func LinkLog() {

	file, err := os.OpenFile(filepath.Join(PgPath, "logs/Pichubot-"+string(time.Now().Format("2006-01-02"))+".log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(filepath.Join(PgPath, "logs/Pichubot-"+string(time.Now().Format("2006-01-02"))+".log"))
			if err != nil {
				panic("[!SERVE!] Could not create the log file. Permission denied.") // maybe
			}
		} else {
			panic("[!SERVE!] Can't read log Please check your permission.")
		}
	}
	LogFile = file

}

// 在屏幕上输出日志并储存
// 日志等级 - int
// DEBUG   -  1
// INFO    -  2
// WARNING -  3
// SEVERE  -  4
func PrintLog(level int, message string) {
	time := time.Now().Format("15:04:05")
	// SEVERE, WARNING, INFO and DEBUG
	switch level {
	// debug level
	case 1:
		if Config.Loglvl == 4 {
			fmt.Printf("[%s][Debug]:%s\n", time, message)
			_, err := LogFile.WriteString(fmt.Sprintf("[%s][Debug] %s\n", time, message))
			if err != nil {
				panic(err)
			}
		}

	// info level
	case 2:
		if Config.Loglvl >= 3 {
			fmt.Printf("[%s][INFO]:%s\n", time, message)
			_, err := LogFile.WriteString(fmt.Sprintf("[%s][INFO] %s\n", time, message))
			if err != nil {
				panic(err)
			}
		}

	// warning level
	case 3:
		if Config.Loglvl >= 2 {
			fmt.Printf("[%s][WARN]:%s\n", time, message)
			_, err := LogFile.WriteString(fmt.Sprintf("[%s][WARN] %s\n", time, message))
			if err != nil {
				panic(err)
			}
		}

	// severe level
	case 4:
		if Config.Loglvl >= 1 {
			fmt.Printf("[%s][!SERVE!]:%s\n", time, message)
			_, err := LogFile.WriteString(fmt.Sprintf("[%s][!SERVE!] %s\n", time, message))
			if err != nil {
				panic(err)
			}
		}

	// sth wrong with level
	default:
		PrintLog(4, "A code error was found in PrintLog")
	}
}

// bool to string
func BoolToStr(B bool) string {
	if B {
		return "true"
	} else if !B {
		return "false"
	}
	return ""
}