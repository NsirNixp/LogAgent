package etcd

import (
	"LogCollectAgent/model"
	"LogCollectAgent/tailf"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/coreos/etcd/mvcc/mvccpb"

	"github.com/astaxie/beego/logs"
	etcdClient "github.com/coreos/etcd/clientv3"
)

var etcd *model.Etcd

func InitEtcd(etcdAddress string, etcdPrefix string, ips []string) (collectionConf []model.CollectConf, err error) {
	mEtcdClient, err := etcdClient.New(etcdClient.Config{
		// 实际工作中，为了防止ip会变动，一般写域名或公网ip
		Endpoints: []string{"127.0.0.1:2379", "127.0.0.1:2380", "127.0.0.1:2381"},
		// 超时时间
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		fmt.Printf("Fail to connect to etcd ,[error=%v]\n", err)
		return
	}

	logs.Debug("Connect to etcd successfully")
	fmt.Printf("Connect to etcd successfully\n")

	// 存mEtcdClient
	etcd = &model.Etcd{Client: mEtcdClient}

	// 容错处理
	if !strings.HasSuffix(etcdPrefix, "/") {
		etcdPrefix = etcdPrefix + "/"
	}

	isNeedError := true
	for _, ip := range ips {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		etcdKey := fmt.Sprintf("%s%s", etcdPrefix, ip)
		etcdResp, errEtcd := mEtcdClient.Get(ctx, etcdKey)
		cancel()
		// 报错读取下一个key：兼容window、mac等可能因为虚拟机，无线网卡等原因有多个ip
		if errEtcd != nil {
			continue
		}
		logs.Debug("etcdResp from etcd:%v", etcdResp.Kvs)
		// 有一个网卡成功获取则无需返回错误
		isNeedError = false

		// 添加etcd的key到结构体中
		etcd.Keys = append(etcd.Keys, etcdKey)

		for _, ev := range etcdResp.Kvs {
			fmt.Printf("Load the config success [key=%s],[value=%s]\n", ev.Key, ev.Value)
			logs.Debug("Load the config success [key=%s],[value=%s]\n", ev.Key, ev.Value)

			if etcdKey == string(ev.Key) {
				err = json.Unmarshal([]byte(ev.Value), &collectionConf)
				if err != nil {
					fmt.Printf("Serialize json error ,[error=%v]\n", err)
					return
				}

				logs.Debug("collectionConf", collectionConf)
				fmt.Println("collectionConf", collectionConf)
			}
		}
	}

	if isNeedError {
		err = fmt.Errorf("Can't Load any configuration from etcd.Please check it out")
		return
	}

	logs.Debug("Get configuration from etcd successfully")
	fmt.Printf("Get configuration from etcd successfully\n")

	// 监听etcd发来的通知
	initEtcdWatcher()
	return
}

func initEtcdWatcher() {
	logs.Debug("Start watcher...")
	for _, key := range etcd.Keys {
		logs.Debug("watcher ... etcdKey = %v", key)
		go watcher(key)
	}
}

func watcher(etcdKey string) {

	etcdctl, err := etcdClient.New(etcdClient.Config{
		Endpoints:   []string{"127.0.0.1:2379", "127.0.0.1:2380", "127.0.0.1:2381"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		fmt.Printf("Fail to connect to etcd ,[error=%v]\n", err)
		return
	}

	logs.Debug("Begin to watch [etcdKey=%s]", etcdKey)

	for {
		ctx := context.Background()
		watchChan := etcdctl.Watch(ctx, etcdKey)
		var collectConf []model.CollectConf
		var isUnmarshalRight = true
		for etcdResp := range watchChan {
			for _, ev := range etcdResp.Events {
				// 删除key
				if ev.Type == mvccpb.DELETE {
					logs.Warn("Key have been deleted in etcd,[key=%s]", etcdKey)
					continue
				}

				// 重设value
				if ev.Type == mvccpb.PUT && etcdKey == string(ev.Kv.Key) {
					err := json.Unmarshal(ev.Kv.Value, &collectConf)
					if err != nil {
						// 反序列化失败，打日志and设置标志，不需要更新配置
						logs.Error("Serialize error,[value=%s],[error=%v]", ev.Kv.Value, err)
						fmt.Printf("Serialize error,[value=%s],[error=%v]\n", ev.Kv.Value, err)
						isUnmarshalRight = false
					}
					continue
				}
				fmt.Printf("type=%s,[key=%s],[value=%s]\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				logs.Debug("type=%s,[key=%s],[value=%s]", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}

			// 需要更新配置
			if isUnmarshalRight {
				logs.Debug("Get configuration from etcd, %v", collectConf)
				tailf.UpdateConfig(collectConf)
			}

		}

	}
}
