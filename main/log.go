package main

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/logs"
)

func levelConverter(appLogLevel string) int {
	switch appLogLevel {
	case TRACE:
		return logs.LevelTrace
	case DEBUG:
		return logs.LevelDebug
	case WARN:
		return logs.LevelWarn
	case INFO:
		return logs.LevelInfo
	case ERROR:
		return logs.LevelError
	default:
		return logs.LevelDebug
	}
}

/*
 * 初始化logger
 */
func initLogger() (err error) {

	config := make(map[string]interface{})
	config["filename"] = appConfig.AppLogPath
	config["level"] = levelConverter(appConfig.AppLogLevel)

	types, err := json.Marshal(config)
	if err != nil {
		fmt.Printf("marshal json error. [error=%v]\n", err)
		return
	}

	logs.SetLogger(logs.AdapterFile, string(types))
	logs.Debug("initLogger success,[config=%v]", string(types))

	return

}
