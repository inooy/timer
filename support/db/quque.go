package db

import (
	"strconv"
	"time"
	"timer/core"
	"timer/pkg/log"
)

type TimerQueueImpl struct {
	Name   string
	Runner func(task *core.DelayTask) error
}

func NewDbTimer(name string, runner func(task *core.DelayTask) error) *TimerQueueImpl {
	return &TimerQueueImpl{
		Name:   name,
		Runner: runner,
	}
}

func (queue *TimerQueueImpl) GetName() string {
	return queue.Name
}

func (queue *TimerQueueImpl) Push(task *core.DelayTask) error {
	log.Infof("push new delay task:[taskId=%s,delay=%d]", task.Id, task.Delay)
	// 数据库主键防重
	po := GormDelayTask{
		Topic:      task.Topic,
		Body:       task.Body,
		Retry:      task.Retry,
		RetryCount: task.RetryCount,
		Delay:      task.Delay,
		Status:     int(task.Status),
	}
	if err := Insert(&po); err != nil {
		return err
	}
	task.Id = strconv.FormatInt(po.Id, 10)
	queue.submit(task)
	log.Info("push success:[taskId=%s]", task.Id)
	return nil
}

func (queue *TimerQueueImpl) ReToDelayQueue(task *core.DelayTask) error {
	log.Infof("task reto delay queue:[id=%s]", task.Id)
	queue.submit(task)
	log.Infof("task reto delay queue success:[id=%s]", task.Id)
	return nil
}

func (queue *TimerQueueImpl) submit(task *core.DelayTask) {
	var pre int32
	n, err := strconv.Atoi(task.Id)
	if err != nil {
		log.Errorf("id error:", err.Error())
		return
	}
	next := int32(n)
	for {
		pre = currentId.Load()
		// id 增大才更新
		if next <= pre {
			break
		}
		if currentId.CompareAndSwap(pre, next) {
			break
		}
	}
	now := time.Now().UnixNano() / 1e6
	delay := task.Delay - (now - task.CreateTime)
	log.Info("[id=%s]Continue delay duration：%d", task.Id, delay)
	if delay <= 10 {
		log.Infof("[id=%s]No need for timing, direct trigger：%d", task.Id, delay)
		go execute(queue.Runner, task)
		return
	}

	time.AfterFunc(time.Duration(delay)*time.Millisecond, func() {
		go execute(queue.Runner, task)
	})
}
