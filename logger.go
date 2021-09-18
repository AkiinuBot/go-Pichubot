package Pichubot

import (
	"os"

	go_logger "github.com/phachon/go-logger"
)

const (
	LOGGER_LEVEL_EMERGENCY = iota
	LOGGER_LEVEL_ALERT
	LOGGER_LEVEL_CRITICAL
	LOGGER_LEVEL_ERROR
	LOGGER_LEVEL_WARNING
	LOGGER_LEVEL_NOTICE
	LOGGER_LEVEL_INFO
	LOGGER_LEVEL_DEBUG
)

var Logger *go_logger.Logger

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

func InitLogger(lvl int) {
	Logger = go_logger.NewLogger()
	Logger.Detach("console")
	consoleConfig := &go_logger.ConsoleConfig{
		Color:      true,
		JsonFormat: false,
		Format:     "%timestamp_format% [%level_string%] %body%",
	}
	Logger.Attach("console", lvl, consoleConfig)

	fileConfig := &go_logger.FileConfig{
		Filename:   "logs/Pichubot.log",
		MaxSize:    1024 * 1024,
		MaxLine:    10000,
		DateSlice:  "d",
		JsonFormat: false,
		Format:     "%millisecond_format% [%level_string%] [%file%:%line%] %body%",
	}
	Logger.Attach("file", lvl, fileConfig)

	Logger.SetAsync()
}
