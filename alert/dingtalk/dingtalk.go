package dingtalk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	whurl string
)

func Init(token string) {
	whurl = fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", token)
}

func Alert(msg string) error {
	// 要发送的消息内容
	message := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": msg,
		},
	}

	// 将消息内容转换为JSON格式
	messageData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// 发送HTTP POST请求到钉钉机器人的Webhook地址
	resp, err := http.Post(whurl, "application/json", bytes.NewBuffer(messageData))
	if err != nil {
		return err
	}
	// 关闭响应体
	resp.Body.Close()

	return nil
}
