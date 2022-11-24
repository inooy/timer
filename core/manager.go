package core

import (
	"errors"
	"fmt"
	"timer/pkg/log"
)

var Manager *DelayQueueManager

type DelayQueueManager struct {
	QueueMap map[string]DelayQueue
}

func (d *DelayQueueManager) Registry(queue DelayQueue) {
	d.QueueMap[queue.GetName()] = queue
}

func (d *DelayQueueManager) Push(task *DelayTask) error {
	if queue, ok := d.QueueMap[task.Topic]; ok {
		err := queue.Push(task)
		if err != nil {
			return err
		}
	}
	return errors.New(fmt.Sprintf("the queue [name=%s] is not register", task.Topic))
}

func (d *DelayQueueManager) Dispatch(task *DelayTask) error {
	// 根据不同的topic，放入不同队列中
	if queue, ok := d.QueueMap[task.Topic]; ok {
		if task.Status == Restore || task.Status == Delay {
			err := queue.ReToDelayQueue(task)
			if err != nil {
				return err
			}
		} else {
			log.Infof("the task status = %d, not push to delay queue", task.Status)
		}
	}
	return errors.New(fmt.Sprintf("the queue [name=%s] is not register", task.Topic))
}
