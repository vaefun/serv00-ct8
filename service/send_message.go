package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/arlettebrook/serv00-ct8/configs"
	"github.com/arlettebrook/serv00-ct8/utils"
)

var cfg = configs.Cfg

func closeRespBody(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		Logger.Warnf("Close resp body error: %s", err)
	}
}
func sendPushPlusMsg(msg string) {
	url := "https://www.pushplus.plus/send"
	data := map[string]any{
		"token":   cfg.PushPlusToken,
		"title":   "serv00&ct8 notify",
		"content": msg,
	}
	bytesData, err := json.Marshal(data)
	if err != nil {
		Logger.Errorf("序列化JSON数据失败: %s", err)
		return
	}
	resp, err := http.Post(url, "application/json",
		bytes.NewReader(bytesData))
	if err != nil {
		Logger.Errorf("发送POST请求失败: %s", err)
		return
	}
	defer closeRespBody(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK || err != nil {
		Logger.WithField("resp", string(body)).Warnf(
			"PushPlus消息推送失败: %s", err)
	} else {
		Logger.WithField("resp", string(body)).Info(
			"PushPlus消息推送成功")
	}

}

func sendTelegramBotMsg(msg string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage",
		cfg.TelegramBotToken)

	data := map[string]any{
		"chat_id": cfg.TelegramChatId,
		"text":    msg,
		"reply_markup": map[string]any{
			"inline_keyboard": [][]map[string]any{
				{
					{
						"text": "问题反馈❓",
						"url":  "https://github.com/arlettebrook/serv00-ct8/issues",
					},
				},
			},
		},
	}

	bytesData, err := json.Marshal(data)
	if err != nil {
		Logger.Errorf("Marshal json error: %s", err)
		return
	}

	resp, err := http.Post(url, "application/json",
		bytes.NewReader(bytesData))
	if err != nil {
		Logger.Errorf("Send post error: %s", err)
		return
	}
	defer closeRespBody(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK || err != nil {
		Logger.WithField("resp", string(body)).Warnf(
			"TelegramBot消息推送失败: %s", err)
	} else {
		Logger.Info("TelegramBot消息推送成功")
		Logger.Debugf("resp: %s", string(body))
	}
}

func SendMessage(msg string) {
	var wg sync.WaitGroup

	if utils.IsNotEmptyStr(cfg.PushPlusToken) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sendPushPlusMsg(msg)
		}()
	}
	if utils.IsNotEmptyStr(cfg.TelegramBotToken) && utils.
		IsNotEmptyStr(cfg.TelegramChatId) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sendTelegramBotMsg(msg)
		}()
	}

	wg.Wait()
}
