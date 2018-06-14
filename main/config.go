package main

import (
	"LogCollectAgent/model"
	"fmt"

	"github.com/astaxie/beego/config"
)

const (
	INI  = "ini"
	YAML = "yaml"
	XML  = "xml"
	JSON = "json"
)

const (
	TRACE = "trace"
	DEBUG = "debug"
	WARN  = "warn"
	INFO  = "info"
	ERROR = "error"
)

const DEFAULT_AMOUNT_COLLECT_CONF = 1

var appConfig *model.Config

/*
 * 加载配置文件
 */
func loadConf(confType, filename string) (err error) {
	conf, err := config.NewConfig(confType, filename)
	if err != nil {
		fmt.Printf("Load the configuration error. [error=%v]\n", err)
		return
	}
	appConfig = &model.Config{}

	// log
	appLogLevel, appLogPath := loadLogs(conf)

	// collect
	chanSize, err := loadCollect(conf)
	if err != nil {
		return
	}

	// kafka
	kafkaProducerAddress, err := loadKafkaConf(conf)
	if err != nil {
		return
	}

	appConfig.AppLogLevel = appLogLevel
	appConfig.AppLogPath = appLogPath
	appConfig.ChanSize = chanSize
	appConfig.KafkaConf.ProducerAddress = kafkaProducerAddress

	fmt.Printf("Load conf finished,[AppConfig=%v]\n", appConfig)
	return
}

/*
 * 读取kafka节点
 */
func loadKafkaConf(conf config.Configer) (kafkaProducerAddress string, err error) {
	kafkaProducerAddress = conf.String("kafka::kafka_producer_address")
	if len(kafkaProducerAddress) == 0 {
		fmt.Printf("Failed to Load configuration [kafka::kafka_producer_address]\n")
		return
	}
	return
}

/*
 * 读取collect节点
 */
func loadCollect(conf config.Configer) (chanSize int, err error) {

	logPath := conf.String("collect::log_path")
	if len(logPath) == 0 {
		fmt.Printf("Failed to Load configuration [collect::log_path]\n")
		return
	}

	topic := conf.String("collect::topic")
	if len(topic) == 0 {
		fmt.Printf("Failed to Load configuration [collect::log_path]\n")
		return
	}

	chanSize, err = conf.Int("collect::chan_size")
	if err != nil {
		chanSize = 100
		return
	}

	collectConf := model.CollectConf{
		LogPath: logPath,
		Topic:   topic,
	}

	collectConfs := append(appConfig.CollectConf, collectConf)
	appConfig.CollectConf = collectConfs
	return
}

/*
 * 读取logs节点
 */
func loadLogs(conf config.Configer) (logLevel, logPath string) {

	logLevel = conf.String("logs::log_level")
	if len(logLevel) == 0 {
		logLevel = DEBUG
	}

	logPath = conf.String("logs::log_path")
	if len(logLevel) == 0 {
		logLevel = "./logs"
	}

	return
}
