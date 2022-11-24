# timer
golang定时任务组件

基于原生go time组件实现定时任务，通过gorm进行数据持久化

## 🎉 特性

1. [x] 原生time组件实现定时任务
2. [x] 支持延迟任务持久化
3. [x] 支持定时扫描数据库是否有新任务，在时间不是特别敏感的情况可以作为分布式系统的兜底策略
4. [x] 支持设置任务消费失败重试

## 💯 使用
### 引入依赖：
```shell
go get github.com/inooy/timer
```

### 更新依赖
```shell
go get github.com/inooy/timer@v0.1.1
```

### 编程使用

```go

package main

import (
	"gorm.io/gorm"
	"timer/core"
	"timer/support/db"
	"timer/timer"
)

func main() {
	var orm gorm.DB
	timer.SetUp(&orm)
	queue := db.NewDbTimer("delay-notify", func(task *core.DelayTask) error {
		return nil
	})
	core.Manager.Registry(queue)

}


```


