package Pichubot

import (
	"encoding/json"
	"os"

	"github.com/wonderivan/logger"
)

/*
等级	配置	释义	控制台颜色
0	EMER	系统级紧急，比如磁盘出错，内存异常，网络不可用等	红色底
1	ALRT	系统级警告，比如数据库访问异常，配置文件出错等	紫色
2	CRIT	系统级危险，比如权限出错，访问异常等	蓝色
3	EROR	用户级错误	红色
4	WARN	用户级警告	黄色
5	INFO	用户级重要	天蓝色
6	DEBG	用户级调试	绿色
7	TRAC	用户级基本输出	绿色
*/

type logconfig struct {
	TimeFormat string `json:"TimeFormat"`
	Console    struct {
		Level string `json:"level"`
		Color bool   `json:"color"`
	}
	File struct {
		Filename string `json:"filename"`
		Level    string `json:"level"`
		Daily    bool   `json:"daily"`
		Maxlines int    `json:"maxlines"`
		Maxsize  int    `json:"maxsize"`
		Maxdays  int    `json:"maxdays"`
		Append   bool   `json:"append"`
		Permit   string `json:"permit"`
	}
	// Conn struct {
	// 	Net            string `json:"net"`
	// 	Addr           string `json:"addr"`
	// 	Level          string `json:"level"`
	// 	Reconnect      bool   `json:"reconnect"`
	// 	ReconnectOnMsg bool   `json:"reconnectOnMsg"`
	// }
}

// 判断文件夹是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func init() {
	if exist, err := pathExists("./logs"); !exist && err == nil {
		err := os.Mkdir("./logs", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

func InitLogger(lvl string) {
	var config logconfig
	config.TimeFormat = "2006-01-02 15:04:05"  // 输出日志开头时间格式
	config.Console.Level = "TRAC"              // 控制台日志输出等级
	config.Console.Color = true                // 控制台日志颜色开关
	config.File.Filename = "logs/Pichubot.log" // 初始日志文件名
	config.File.Level = "TRAC"                 // 日志文件日志输出等级
	config.File.Daily = true                   // 跨天后是否创建新日志文件，当append=true时有效
	config.File.Maxlines = 1000000             // 日志文件最大行数，当append=true时有效
	config.File.Maxsize = 1                    // 日志文件最大大小，当append=true时有效
	config.File.Maxlines = -1                  // 日志文件有效期
	config.File.Append = true                  // 是否支持日志追加
	config.File.Permit = "0660"                // 新创建的日志文件权限属性
	data, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	logger.SetLogger(string(data[:]))
}

// "TimeFormat":"2006-01-02 15:04:05",
// "Console": {            // 控制台日志配置
// 		"level": ,
// 		"color": true
// },
// "File": {                   // 文件日志配置
// 		"filename": "app.log",
// 		"level": "TRAC",        // 日志文件日志输出等级
// 		"daily": true,          // 跨天后是否创建新日志文件，当append=true时有效
// 		"maxlines": 1000000,    // 日志文件最大行数，当append=true时有效
// 		"maxsize": 1,           // 日志文件最大大小，当append=true时有效
// 		"maxdays": -1,          // 日志文件有效期
// 		"append": true,         // 是否支持日志追加
// 		"permit": "0660"        // 新创建的日志文件权限属性
// },
// "Conn": {                       // 网络日志配置
// 		"net":"tcp",                // 日志传输模式
// 		"addr":"10.1.55.10:1024",   // 日志接收服务器
// 		"level": "Warn",            // 网络日志输出等级
// 		"reconnect":true,           // 网络断开后是否重连
// 		"reconnectOnMsg":false,     // 发送完每条消息后是否断开网络
