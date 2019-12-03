# dingding
golang实现钉钉异步报警,可限制报警次数,具备超时机制
## usage

```shell
    go get github.com/wulorn/dingding
```

```go
    // 主程或初始化时调用即可
    go dingding.WithRecovery(dingding.SendAlarmMessAge)
    // 发送告警信息的地方调用即可
    dingding.PushMessage("alarm xxx")
    
	select {}
```