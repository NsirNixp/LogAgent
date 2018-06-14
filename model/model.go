package model

import (
	"github.com/hpcloud/tail"
)

/**
 * LogAgent系统的配置对象
 */
type Config struct {
	AppLogLevel string        // LogAgent日志等级
	AppLogPath  string        // LogAgent日志存放地址
	CollectConf []CollectConf // 收集日志信息基本
	ChanSize    int
	KafkaConf
}

/**
 * 一台机器的收集配置信息
 */
type CollectConf struct {
	LogPath string // 收集日志存放地址
	Topic   string // kafka消费topic
}

type KafkaConf struct {
	ProducerAddress string
}

/**
 * Kafka消息
 */
type TextMsg struct {
	Msg   string
	Topic string
}

/**
 * 日志扫描信息
 */
type TailObject struct {
	Tail *tail.Tail  // 日志查看对象
	Conf CollectConf // 一台机器的收集配置
}

/**
 * 日志扫描信息列表
 */
type TailList struct {
	TailList []*TailObject
	MsgChan  chan *TextMsg // 发送kafka消息的channel
}
