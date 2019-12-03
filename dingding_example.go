package dingding

import (
	"fmt"
	"log"
	"time"
)

const (
	// 一般通过配置文件配置
	openUrl = "https://oapi.dingtalk.com/robot/send?access_token=xxx" //"钉钉生成的url"
	mobiles = "18012345678,18212345678"                               // @ 通知对象手机号
	MaxNum  = 1<<32 - 1
)

func SendAlarmMessAge() {
	ch := GetDingCh()
	alarmCntTimer := time.NewTimer(time.Second * 60)
	var alarmCnt, lastAlarmCnt int
	for {
		select {
		case ctx := <-ch:
			alarmCnt++
			if (alarmCnt - lastAlarmCnt) <= 10 { // 每分钟最大允许发送告警次数
				alarmStr := fmt.Sprintf("时间: %s 服务xxx 告警信息: %s", time.Now().Format("2006-01-02 15:04:05"), ctx)
				SendDingCh(openUrl, mobiles, alarmStr)
			}
		case <-alarmCntTimer.C:
			if alarmCnt >= MaxNum {
				alarmCnt = 0
			}
			lastAlarmCnt = alarmCnt
			log.Printf("%d alarms have been sent.", alarmCnt)
			alarmCntTimer.Reset(time.Second * 60)
		}
	}
}

var PanicHandler func(interface{})

func WithRecovery(f func()) {
	defer func() {
		if err := recover(); err != nil {
			if PanicHandler != nil {
				PanicHandler(err)
			}
		}
	}()
	f()
}

func main() {
	// usage
	go WithRecovery(SendAlarmMessAge)
	PushMessage("alarm xxx") // 发送告警信息
	select {}
}
