package model

import (
	"sync"

	etcdClient "github.com/coreos/etcd/clientv3"
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
	EtcdConf
}

/**
 * 一台机器的收集配置信息
 */
type CollectConf struct {
	LogPath string `json:"logpath"` // 收集日志存放地址
	Topic   string `json:"topic"`   // kafka消费topic
}

type KafkaConf struct {
	Address string
}

type EtcdConf struct {
	Address string
	Key     string
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
	Tail     *tail.Tail  // 日志查看对象
	Conf     CollectConf // 一台机器的收集配置
	Status   int         // 状态：标记当前tail实例状态
	ExitChan chan bool   // 退出标记：true == <- ExitChan 表示该channel已完成/被结束
}

/**
 * 日志扫描信息列表
 */
type TailList struct {
	TailList []*TailObject
	MsgChan  chan *TextMsg // 发送kafka消息的channel
	// Lock     sync.RWMutex
	Lock sync.Mutex
}

/**
 * etcd实例
 */
type Etcd struct {
	Client *etcdClient.Client // etcd客户端实例
	Keys   []string           // 本机日志文件配置对应在etcd中的key {"/logagent/conf/logs/178.156.133.22","/logagent/conf/logs/178.156.133.21"}
}
