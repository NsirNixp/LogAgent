package main

/**
 * @desc  : 基础数据发送到etcd服务器，便于测试
 * @author: Charelyz
 */

import (
	"LogCollectAgent/model"
	"context"
	"encoding/json"
	"fmt"
	"time"

	etcdClient "github.com/coreos/etcd/clientv3"
)

const (
	DEFAULT_ETCD_KEY = "/logagent/conf/logs/169.254.237.219"
)

var (
	collectionConf []model.CollectConf
)

func sendToEtcd() {
	etcdctl, err := etcdClient.New(etcdClient.Config{
		Endpoints:   []string{"127.0.0.1:2379", "127.0.0.1:2380", "127.0.0.1:2381"},
		DialTimeout: time.Second * 5,
	})
	defer etcdctl.Close()
	if err != nil {
		fmt.Printf("Fail to connect to etcd ,[error=%v]\n", err)
		return
	}

	// 准备测试数据
	collectionConf = append(collectionConf, model.CollectConf{
		LogPath: "D:/4-学习空间/golang/learning/logs/log_agent.log",
		Topic:   "nginx_access",
	})
	collectionConf = append(collectionConf, model.CollectConf{
		LogPath: "D:/4-学习空间/golang/learning/logs/client/error.log",
		Topic:   "nginx_error",
	})

	types, err := json.Marshal(collectionConf)
	if err != nil {
		fmt.Printf("Serialize json error ,[error=%v]\n", err)
		return
	}

	// 操作etcd
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	_, err = etcdctl.Put(ctx, DEFAULT_ETCD_KEY, string(types))
	cancel()
	if err != nil {
		fmt.Printf("Fail to set value in etcd ,[error=%v]\n", err)
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*1)
	etcdResp, err := etcdctl.Get(ctx, DEFAULT_ETCD_KEY)
	cancel()
	if err != nil {
		fmt.Printf("Fail to get value from etcd ,[error=%v]\n", err)
		return
	}

	for _, ev := range etcdResp.Kvs {
		fmt.Printf("key=%s\nvalue=%s\n", ev.Key, ev.Value)
	}
}

func sendToEtcdDelete() (err error) {
	etcdctl, err := etcdClient.New(etcdClient.Config{
		Endpoints:   []string{"127.0.0.1:2379", "127.0.0.1:2380", "127.0.0.1:2381"},
		DialTimeout: time.Second * 5,
	})
	defer etcdctl.Close()
	if err != nil {
		fmt.Printf("Fail to connect to etcd ,[error=%v]\n", err)
		return
	}

	// 操作etcd
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	_, err = etcdctl.Delete(ctx, DEFAULT_ETCD_KEY)
	cancel()
	if err != nil {
		fmt.Printf("Fail to set value in etcd ,[error=%v]\n", err)
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*1)
	etcdResp, err := etcdctl.Get(ctx, DEFAULT_ETCD_KEY)
	cancel()
	if err != nil {
		fmt.Printf("Fail to get value from etcd ,[error=%v]\n", err)
		return
	}

	for _, ev := range etcdResp.Kvs {
		fmt.Printf("key=%q\nvalue=%q\n", DEFAULT_ETCD_KEY, ev.Value)
	}
	fmt.Printf("delete success")
	return
}

func main() {
	// sendToEtcd()
	sendToEtcdDelete()
}
