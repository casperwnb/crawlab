package msg_handler

import (
	"crawlab/constants"
	"crawlab/entity"
	"crawlab/model"
	"crawlab/utils"
	"github.com/apex/log"
	"runtime/debug"
)

type Log struct {
	msg entity.NodeMessage
}

func (g *Log) Handle() error {
	if g.msg.Type == constants.MsgTypeGetLog {
		return g.get()
	} else if g.msg.Type == constants.MsgTypeRemoveLog {
		return g.remove()
	}
	return nil
}

func (g *Log) get() error {
	// 发出的消息
	msgSd := entity.NodeMessage{
		Type:   constants.MsgTypeGetLog,
		TaskId: g.msg.TaskId,
	}
	// 获取本地日志
	logStr, err := model.GetLocalLog(g.msg.LogPath)
	if err != nil {
		log.Errorf("get node local log error: %s", err.Error())
		debug.PrintStack()
		msgSd.Error = err.Error()
		msgSd.Log = err.Error()
	} else {
		msgSd.Log = utils.BytesToString(logStr)
	}
	// 发布消息给主节点
	if err := utils.Pub(constants.ChannelMasterNode, msgSd); err != nil {
		return err
	}
	return nil
}

func (g *Log) remove() error {
	return model.RemoveFile(g.msg.LogPath)
}
