package pichumod

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hjson/hjson-go"
)

// This package is for settings.hjson

var PgLoc, _ = os.Executable()
var PgPath = filepath.Dir(PgLoc)

var Config ConfigStu

func ReadSettings() {
	file, err := os.Open("settings.hjson")
	if err != nil {
		if os.IsNotExist(err) {
			createSettings()
			ReadSettings()
		} else {
			panic("[!SERVE!] Can't read settings.hjson! Please check your permission.")
		}
	} else {
		content, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
		var data map[string]interface{}
		if err := hjson.Unmarshal(content, &data); err != nil {
			panic(err)
		}
		ConfigParse(data)
	}
}

// 解析配置 Config
func ConfigParse(data map[string]interface{}) {
	Config.Logmode = data["logmode"].(string)
	Config.WSIP = data["ws-ip"].(string)
	switch Config.Logmode {
	case "NONE":
		Config.Loglvl = 0
	case "DEBUG":
		Config.Loglvl = 4
	case "INFO":
		Config.Loglvl = 3
	case "WARN":
		Config.Loglvl = 2
	case "SERVE":
		Config.Loglvl = 1
	default:
		panic("Found something wrong in settings.hjson.logmode")
	}
}

func createSettings() {
	setting := []byte(`{
	// 输出的日志等级 DEBUG INFO WARN SERVE NONE
	logmode: INFO

	// cqhttp 的websocket服务器地址
	ws-ip: "localhost:6700"
}`)
	settingFile, err := os.Create(filepath.Join(PgPath, "settings.hjson"))
	if err != nil {
		panic(err)
	}
	settingFile.Write(setting)
	fmt.Println("默认配置文件已经生成 请编辑 settings.hjson 后重启程序")
	settingFile.Close()
}
