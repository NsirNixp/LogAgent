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

/*
 * 初始化taif读取配置文件
 */
func InitTailf(collectConf []model.CollectConf, chanSize int) (err error) {

	if len(collectConf) == 0 {
		logs.Error("Section collect is empty,please check out the [collect] section")
		err = fmt.Errorf("Section collect is empty,please check out the [collect] section")
		return
	}

	tailList = &model.TailList{
		MsgChan: make(chan *model.TextMsg, chanSize),
	}

	for _, cc := range collectConf {

		tailObject := &model.TailObject{
			Conf: cc,
		}

		// 创建滚动文件
		tails, errTail := tail.TailFile(cc.LogPath, tail.Config{
			ReOpen: true, // 重新打开重建的文件(tail -F)
			Follow: true, // 继续读取新行(tail -f)
			// Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 记录上次读取的位置 offset:0 代表不偏移，whence:2
			MustExist: false, // 文件不存在，不会panic，而是等待
			Poll:      true,  // 轮询文件改变
		})

		if errTail != nil {
			err = errTail
			logs.Error("Tailing file error,Please check out the configuration.[error=%v]\n", err)
			err = fmt.Errorf("Tailing file error,Please check out the configuration.[error=%v]", err)
			return
		}

		tailObject.Tail = tails
		tailList.TailList = append(tailList.TailList, tailObject)

		// 启动一个goroutine处理日志
		go startUp(tailObject)

	}

	logs.Debug("InitTailf success,[collectConf=%v]", collectConf)
	return
}

func startUp(tailObject *model.TailObject) {
	// TODO 改为flag结束
	for true {
		line, ok := <-tailObject.Tail.Lines
		if !ok {
			logs.Warn("tail file close reopen, filename:%s\n", tailObject.Tail.Filename)
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

	}
}

func FetchOneMsgFromChannel() (textMsg *model.TextMsg) {
	textMsg = <-tailList.MsgChan
	return
}
