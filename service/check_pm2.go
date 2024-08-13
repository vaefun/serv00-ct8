package service

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/arlettebrook/serv00-ct8/models"
	"github.com/arlettebrook/serv00-ct8/utils"
)

// CheckPM2 检查pm2进程是否运行，没有重新启动pm2进程
func CheckPM2() {
	Logger.Info("Start check pm2...")
	var wg sync.WaitGroup
	msgChan := make(chan string, 3)
	msg := ""
	var total, recovery atomic.Int32
	for _, u := range cfg.Accounts {
		wg.Add(1)
		go check(u, &wg, msgChan, &total, &recovery)
	}

	go func() {
		wg.Wait()
		close(msgChan)
	}()

	for v := range msgChan {
		msg += v
	}
	if recovery.Load() != 0 {
		msg = fmt.Sprintf("检测到pm2未正常运行，一共%d个任务，需要恢复%d个任务"+
			"。所有结果如下：\n\n", total.Load(), recovery.Load()) + msg
		SendMessage(msg)
	}
	Logger.Info("End check pm2")

}

func check(u models.Account, wg *sync.WaitGroup, msgChan chan<- string, total,
	recovery *atomic.Int32) {
	defer wg.Done()

	if !u.IsCheck {
		return
	}

	total.Add(1)
	withField := Logger.WithField(u.Addr, u.Username)

	client, err := ConnSSH(u)
	if err != nil {
		withField.Errorf("Connection ssh error: %s", err)
		msgChan <- fmt.Sprintf("%s: %s 登录失败(%s) %s \n\n", u.Addr,
			u.Username, utils.NowCST(), err)
		return
	}
	defer CloseClient(client)

	session, err := CreateSession(client)
	if err != nil {
		withField.Errorf("Create session error %s", err)
		msgChan <- fmt.Sprintf("%s: %s 登录失败(%s) %s \n\n", u.Addr,
			u.Username, utils.NowCST(), err)
		return
	}
	defer CloseSession(session)

	assertRunPM2Err := func(err error, status []byte) bool {
		if err != nil {
			withField.Errorf("Run pm2 status error: %s", err)
			withField.Debug(string(status))
			msgChan <- fmt.Sprintf("%s: %s 运行pm2失败(%s) %s \n\n", u.Addr,
				u.Username, utils.NowCST(), err)
			return false
		}
		return true
	}

	status, err := session.CombinedOutput(fmt.Sprintf(
		"~/.npm-global/bin/pm2 status"))
	if !assertRunPM2Err(err, status) {
		return
	}

	if bytes.Contains(status, []byte("online")) {
		withField.Info("无需恢复pm2...")
		withField.Debug(string(status))
		msgChan <- fmt.Sprintf("%s: %s 无需恢复pm2(%s) \n\n", u.Addr,
			u.Username, utils.NowCST())
		return
	}

	recovery.Add(1)
	withField.Warnf("检测到pm2进程未正常运行，正在恢复...")
	withField.Debugf("未恢复前: %s", string(status))

	session2, err := CreateSession(client)
	if err != nil {
		Logger.Errorf("Create session error: %s", err)
		msgChan <- fmt.Sprintf("%s: %s 登录失败(%s) %s \n\n", u.Addr,
			u.Username, utils.NowCST(), err)
		return
	}
	defer CloseSession(session2)

	status2, err := session2.CombinedOutput(fmt.Sprintf(
		"~/.npm-global/bin/pm2 resurrect"))
	if !assertRunPM2Err(err, status2) {
		return
	}

	msgChan <- fmt.Sprintf("%s: %s pm2未正常运行(%s) 恢复结果如下:\n %s "+
		"\n\n", u.Addr, u.Username, utils.NowCST(), string(status2))

	withField.Debugf("恢复后: %s", string(status2))
}
