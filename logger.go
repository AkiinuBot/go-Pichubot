package pichubot

import (
	go_logger "github.com/phachon/go-logger"
)

const (
	LOGGER_LEVEL_EMERGENCY = iota // 系统级紧急，比如磁盘出错，内存异常，网络不可用等
	LOGGER_LEVEL_ALERT            // 系统级警告，比如数据库访问异常，配置文件出错等
	LOGGER_LEVEL_CRITICAL         // 系统级危险，比如权限出错，访问异常等
	LOGGER_LEVEL_ERROR            // 用户级错误
	LOGGER_LEVEL_WARNING          // 用户级警告
	LOGGER_LEVEL_NOTICE           // 用户级重要
	LOGGER_LEVEL_INFO             // 用户级提示
	LOGGER_LEVEL_DEBUG            // 用户级调试
)

// 日志接口
var Logger *go_logger.Logger

func init() {
	CheckPath("./logs")
	Logger = new(go_logger.Logger)
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
