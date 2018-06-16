package tailf

import (
	"LogCollectAgent/model"
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
)

// 包含所有的tail实例，以及各个实例对应的信息(log位置，channel的大小，kafka的topic)
var tailList *model.TailList

const (
	RUNNING  = iota // 0
	STOPPING        // 1
)

/*
 * 初始化taif读取配置文件
 */
func InitTailf(collectConf []model.CollectConf, chanSize int) (err error) {

	tailList = &model.TailList{
		MsgChan: make(chan *model.TextMsg, chanSize),
	}

	if len(collectConf) == 0 {
		//logs.Error("Section collect is empty,please check out the [collect] section")
		//err = fmt.Errorf("Section collect is empty,please check out the [collect] section")
		logs.Error("Can't read config data from etcd.Are you forget it?")
		return
	}

	for _, cc := range collectConf {
		newTailfTask(cc)
	}

	logs.Debug("InitTailf success,[collectConf=%v]", collectConf)
	return
}

/**
 * 读取配置文件，并发送到MsgChan等待被接收
 */
func startUp(tailObject *model.TailObject) {
	// TODO 改为flag结束
	for {

		select {
		// goroutine主动退出
		case <-tailObject.ExitChan:
			logs.Warn("Tailf instance exit, [conf=%v]", tailObject.Conf)
			return
		// 读取数据
		case line, ok := <-tailObject.Tail.Lines:
			if !ok {
				logs.Warn("tail file close reopen, filename:%s\n,retry after 100ms", tailObject.Tail.Filename)
				// wait 100ms
				time.Sleep(time.Millisecond * 100)
				continue
			}
			textMsg := &model.TextMsg{
				Msg:   line.Text,
				Topic: tailObject.Conf.Topic,
			}
			// 发送channel
			tailList.MsgChan <- textMsg
			fmt.Println(textMsg)
		}

	}
}

/**
 * 从channel中获取数据
 */
func FetchOneMsgFromChannel() (textMsg *model.TextMsg) {
	textMsg = <-tailList.MsgChan
	return
}

/**
 * 更新配置文件
 */
func UpdateConfig(collectConf []model.CollectConf) {

	tailList.Lock.Lock()
	defer tailList.Lock.Unlock()

	// 待比较的配置列表：通过比较"正在运行的配置列表"和"待比较的配置列表"对比，确定需要操作的tailf实例
	for _, cc := range collectConf {
		var isRunning = false
		// 遍历正在运行的配置
		for _, onlineCC := range tailList.TailList {
			// 配置相同，剔除(不做处理)
			if cc.LogPath == onlineCC.Conf.LogPath {
				isRunning = true
				break
			}
		}

		// 任务在跑，不做处理
		if isRunning {
			continue
		}

		newTailfTask(cc)

	}

	// 保存当前在线的实例
	var tempOnlineTailList []*model.TailObject
	// 刷选出删除状态的tailf实例
	for _, onlineCC := range tailList.TailList {
		onlineCC.Status = STOPPING
		for _, cc := range collectConf {
			if cc.LogPath == onlineCC.Conf.LogPath {
				onlineCC.Status = RUNNING
				break
			}
		}

		// 说明tailf实例被删除了(任务发生变化，取消该实例)
		if onlineCC.Status == STOPPING {
			onlineCC.ExitChan <- true
			continue
		}

		tempOnlineTailList = append(tempOnlineTailList, onlineCC)
	}

	// 重新赋值
	tailList.TailList = tempOnlineTailList
	fmt.Println("tailList.TailList", tailList.TailList)
}

/**
 * 启动tailf任务
 */
func newTailfTask(cc model.CollectConf) {
	tailObject := &model.TailObject{
		Conf:     cc,
		ExitChan: make(chan bool, 1),
	}

	// 创建滚动文件实例
	tails, err := tail.TailFile(cc.LogPath, tail.Config{
		ReOpen: true, // 重新打开重建的文件(tail -F)
		Follow: true, // 继续读取新行(tail -f)
		// Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 记录上次读取的位置 offset:0 代表不偏移，whence:2
		MustExist: false, // 文件不存在，不会panic，而是等待
		Poll:      true,  // 轮询文件改变
	})

	if err != nil {
		logs.Error("Tailing file error [file=%s],Please check out the configuration.[error=%v]\n", cc.LogPath, err)
		fmt.Printf("Tailing file error [file=%s],Please check out the configuration.[error=%v]", cc.LogPath, err)
	}

	tailObject.Tail = tails
	tailList.TailList = append(tailList.TailList, tailObject)

	// 启动一个goroutine处理日志
	go startUp(tailObject)
}
