package db

import (
	"strconv"
	"sync"
	"time"
	"timer/core"
	"timer/pkg/log"
)

const (
	DefaultAccuracy      int64 = 1000
	DefaultCleanInterval int64 = 3 * 24
	DefaultInterval      int64 = 60
	DefaultTimeout       int64 = 60000
)

type QueuesManager struct {
	NeedStop      bool
	Max           int
	Accuracy      int64
	StartTime     int64
	executeTimes  int64
	lastCleanTime int64
	cleanInterval int64
	// 线程扫描间隔 单位秒
	Interval int64
	Timeout  int64
	once     sync.Once
}

type Options struct {
	CleanInterval int64
	Accuracy      int64
	Interval      int64
	Timeout       int64
}

func NewManager(options *Options) *QueuesManager {
	now := time.Now().UnixNano() / 1e6
	if options.CleanInterval == 0 {
		options.CleanInterval = DefaultCleanInterval
	}
	if options.Accuracy == 0 {
		options.Accuracy = DefaultAccuracy
	}
	if options.Accuracy == 0 {
		options.Interval = DefaultInterval
	}
	if options.Timeout == 0 {
		options.Timeout = DefaultTimeout
	}
	return &QueuesManager{
		StartTime:     now,
		cleanInterval: options.CleanInterval,
		Accuracy:      options.Accuracy,
		Interval:      options.Interval,
		Timeout:       options.Timeout,
	}
}

// Launch 启动，循环扫描表
func (m *QueuesManager) Launch() {
	m.once.Do(func() {
		m.run()
	})
}

// Cancel 不会立即停止线程，而是标记需要停止，等待本次轮询执行完自然退出
// 一般在1秒钟内结束
func (m *QueuesManager) Cancel() {
	m.NeedStop = true
}

// IsCanceled 判断调度是否已经停止
// true：调度已取消 false: 调度未取消
func (m *QueuesManager) IsCanceled() bool {
	return m.NeedStop
}

func (m *QueuesManager) run() {
	for {
		if m.NeedStop {
			return
		}
		m.mainLoop()
	}
}

func (m *QueuesManager) mainLoop() {
	log.Info("轮询线程扫描")
	m.executeTimes++
	m.scanNew()
	m.scanRestore()
	m.scanTimeout()
	m.handleClean()
	m.sleep()
}

func (m *QueuesManager) scanNew() {
	// 查询是否有新增：用于同步集群其他实例新增
	list, err := QueryNew(int(currentId.Load()))
	if err != nil {
		log.Error("scan new task error.", err.Error())
	}
	if len(list) == 0 {
		return
	}

	log.Infof("find new task:[count=%d]", len(list))
	dispatch(list)
}

func (m *QueuesManager) scanRestore() {
	list, err := QueryByStatus(core.Restore)
	if err != nil {
		log.Error("scan restore task error.", err.Error())
	}
	if len(list) == 0 {
		return
	}

	log.Infof("find restore task:[count=%d]", len(list))
	dispatch(list)
}
func (m *QueuesManager) scanTimeout() {
	now := time.Now().UnixNano() / 1e6
	list, err := QueryTimeout(now - m.Timeout)
	if err != nil {
		log.Error("scan timeout task error.", err.Error())
	}
	if len(list) == 0 {
		return
	}
	log.Infof("find timout task:[count={}]", len(list))
	dispatch(list)
}

func dispatch(list []GormDelayTask) {
	for i := range list {
		// 还需要检测重试队列、超时队列，死亡队列需要告警人工处理
		task := core.DelayTask{
			Id:         strconv.FormatInt(list[i].Id, 10),
			Topic:      list[i].Topic,
			Body:       list[i].Body,
			Retry:      list[i].Retry,
			RetryCount: list[i].RetryCount,
			Delay:      list[i].Delay,
			Status:     core.TaskStatus(list[i].Status),
			CreateTime: list[i].CreateTime,
			Version:    list[i].Version,
		}
		err := core.Manager.Dispatch(&task)
		if err != nil {
			log.Error("dispatch delay task err.", err.Error())
			return
		}
	}
}

func (m *QueuesManager) handleClean() {
	now := time.Now().UnixNano() / 1e6
	if now-m.lastCleanTime < m.cleanInterval*60*60*1000 {
		return
	}

	if now/int64(60*60*1000)%m.cleanInterval == int64(0) {
		log.Info("start clean task history")
		count, err := CleanOld(now - (m.cleanInterval * 60 * 60 * 1000))
		if err != nil {
			log.Error("clean task history error.", err.Error())
			return
		}
		log.Infof("clean history success，[count=%d]", count)
		m.lastCleanTime = now
	}
}

func (m *QueuesManager) sleep() {
	time.Sleep(time.Second * time.Duration(m.Interval))
}
