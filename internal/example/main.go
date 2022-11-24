package main

import (
	"fmt"
	"gorm.io/gorm"
	"time"
	"timer/core"
	"timer/support/db"
	"timer/timer"
)

func main() {
	var orm gorm.DB
	timer.SetUp(&orm)
	queue := db.NewDbTimer("delay-notify", func(task *core.DelayTask) error {
		fmt.Println("触发延迟任务：", task.Body)
		return nil
	})
	core.Manager.Registry(queue)

	now := time.Now().UnixNano() / 1e6

	task := core.DelayTask{
		Topic:      "delay-notify",
		Body:       "hello im is delay task",
		Retry:      3,
		Delay:      10000,
		CreateTime: now,
	}
	err := core.Manager.Push(&task)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

}
