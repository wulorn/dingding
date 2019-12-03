package dingding

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

var dingCh chan string

const (
	dingLen     = 10 * 1000
	dingTimeOut = time.Second * 1
)

func init() {
	dingCh = make(chan string, dingLen)
}

func PushMessage(msg string) {
	if len(dingCh) >= cap(dingCh) {
		log.Println("error: push message failed. dingCh is full.")
		return
	}
	dingCh <- msg
}

func GetDingCh() chan string {
	return dingCh
}

// openUrl 钉钉的链接、 mobiles: 消息通知对象手机号(字符串 以逗号分割)、 ctx: 具体发送的文本消息
func SendDingCh(openUrl, mobiles, ctx string) {
	postData := make(map[string]interface{})
	postData["msgtype"] = "text"
	sendContext := make(map[string]interface{})
	sendContext["content"] = ctx
	postData["text"] = sendContext

	at := make(map[string]interface{})
	at["atMobiles"] = strings.Split(mobiles, ",") // 根据手机号@指定人
	at["isAtAll"] = false                         // 禁用@所有人
	postData["at"] = at

	msg, err := json.Marshal(postData)
	if err != nil {
		log.Println("ding ding postData json marshal err => ", err)
		return
	}

	// 通过http.Client 中的 DialContext 可以设置连接超时和数据接受超时 （也可以使用Dial, 不推荐）
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				conn, err := net.DialTimeout(network, addr, dingTimeOut) // 设置建立链接超时
				if err != nil {
					return nil, err
				}
				_ = conn.SetDeadline(time.Now().Add(dingTimeOut)) // 设置接受数据超时时间
				return conn, nil
			},
			ResponseHeaderTimeout: dingTimeOut, // 设置服务器响应超时时间
		},
	}
	req, err := http.NewRequest("POST", openUrl, strings.NewReader(string(msg)))
	if err != nil {
		log.Println("ding ding new post request err =>", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Println("ding ding post request err =>", err)
		return
	}
	defer resp.Body.Close()
}
