package db

import (
	"strconv"
	"timer/core"
	"timer/pkg/log"
)

func execute(runner func(task *core.DelayTask) error, task *core.DelayTask) {
	log.Infof("定时消息到期触发：[id=%s, body=%+v]", task.Id, task.Body)
	id, err := strconv.Atoi(task.Id)
	if err != nil {
		log.Errorf("id不是number,[id=%s]", task.Id)
		return
	}
	dbMessage, err := QueryById(id)
	if err != nil {
		log.Errorf("消息数据丢失,[id=%s]", task.Id)
		return
	}
	// 如果正在执行中，那么不管 -> 对于执行超时的，单独去扫描，加入重试队列中去，这里直接判定为重复执行
	if dbMessage.Status != int(core.Delay) {
		log.Warnf("消息状态不正确，可能有其他服务已处理,[id=%s, status=%d]", task.Id, dbMessage.Status)
		return
	}
	// 放入处理中
	result := pushToProcessing(dbMessage)
	if result != 1 {
		log.Warnf("修改消息为processing失败，可能其他服务正在处理,[id=%s]", task.Id)
		return
	}
	if err = runner(task); err != nil {
		log.Errorf("消费消息失败,[id=%s]", task.Id)
		if dbMessage.RetryCount >= dbMessage.Retry {
			log.Errorf("重试次数达到限制,[id=%s, retry=%d]", task.Id, dbMessage.RetryCount)
			pushToDeadQueue(dbMessage)
		} else {
			log.Warnf("加入重试队列,[id=%s, retry=%d]", task.Id, dbMessage.RetryCount)
			pushToRestoreQueue(dbMessage)
		}
	}
	pushToSuccess(dbMessage)
	log.Infof("定时消息消费成功：[id=%s]", task.Id)
}

func pushToProcessing(po *GormDelayTask) int64 {
	result, err := UpdateStatus(po.Id, core.Processing, po.Version)
	if err != nil {
		log.Warnf("修改消息为processing失败，可能其他服务正在处理,[id=%d]", po.Id)
		return 0
	}
	po.Version++
	return result
}

func pushToSuccess(po *GormDelayTask) int64 {
	result, err := UpdateStatus(po.Id, core.Success, po.Version)
	if err != nil {
		log.Warnf("修改消息为success失败，可能其他服务正在处理,[id=%d]", po.Id)
		return 0
	}
	po.Version++
	return result
}
func pushToRestoreQueue(po *GormDelayTask) int64 {
	result, err := UpdateStatus(po.Id, core.Restore, po.Version)
	if err != nil {
		log.Warnf("修改消息为restore失败，可能其他服务正在处理,[id=%d]", po.Id)
		return 0
	}
	po.Version++
	return result
}

func pushToDeadQueue(po *GormDelayTask) int64 {
	result, err := UpdateStatus(po.Id, core.Dead, po.Version)
	if err != nil {
		log.Warnf("修改消息为dead失败，可能其他服务正在处理,[id=%d]", po.Id)
		return 0
	}
	po.Version++
	return result
}
