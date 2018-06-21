package main

import (
	"LogCollectAgent/etcd"
	"LogCollectAgent/kafka"
	"LogCollectAgent/tailf"
	"fmt"

	"github.com/astaxie/beego/logs"
)

/*
 * @author: Charelyz
 */

func main() {
	// 加载配置文件
	// filename := "./src/LogCollectAgent/conf/logAgent.conf"
	// filename := "D:/4-学习空间/golang/learning/src/LogCollectAgent/conf/logAgent.conf"
	filename := "./conf/log_agent.conf"
	err := loadConf(INI, filename)
	if err != nil {
		fmt.Printf("Load the configuration error. [error=%v]\n", err)
		panic("Failed to Load configuration.Please make sure that the configuration exists")
	}
	fmt.Printf("loadConf()...\n")

	// 初始化logger
	err = initLogger()
	if err != nil {
		fmt.Printf("initLogger error. [error=%v]\n", err)
		panic("Failed to init logger component.Please make sure that the configuration is right")
	}
	fmt.Printf("initLogger()...\n")

	// go func() {
	// 	var count int
	// 	for {
	// 		count++
	// 		logs.Info("info %d", count)
	// 		time.Sleep(time.Second * 1)
	// 	}
	// }()

	// 初始化etcd
	collectionConf, err := etcd.InitEtcd(appConfig.EtcdConf.Address, appConfig.EtcdConf.Key, ips)
	if err != nil {
		fmt.Printf("InitEtcd error. [error=%v]\n", err)
		panic("Failed to init etcd component.Please make sure that the configuration is right")
	}

	// 初始化tailf
	err = tailf.InitTailf(collectionConf, appConfig.ChanSize)
	if err != nil {
		fmt.Printf("initTailf error. [error=%v]\n", err)
		panic("Failed to init tailf component.Please make sure that the configuration is right")
	}
	fmt.Printf("tailf.InitTailf()...\n")

	// 初始化kafka
	err = kafka.InitKafka(appConfig.KafkaConf.Address)
	if err != nil {
		fmt.Printf("InitKafka error. [error=%v]\n", err)
		panic("Failed to init kafka component.Please make sure that the configuration is right")
	}

	// 执行LogAgent
	err = Run()
	if err != nil {
		fmt.Printf("Server have stopped. [error=%v]\n", err)
		panic("Failed to run.Maybe something wrong, please check the log as soon as possible. ")
	}

	logs.Info("program exited")

}
