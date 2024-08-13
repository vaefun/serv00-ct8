package service

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/arlettebrook/serv00-ct8/models"
	"github.com/arlettebrook/serv00-ct8/utils"
)

// TerminalLogin 登录远程服务器
func TerminalLogin() {
	Logger.Info("Start terminal login...")
	var wg sync.WaitGroup
	msgChan := make(chan string, 3)
	var countOpt atomic.Int32
	msg := "serv00&ct8 ssh自动登录结果如下: \n\n"
	for _, u := range cfg.Accounts {
		wg.Add(1)
		go login(u, &wg, msgChan, &countOpt)
	}
	go func() {
		wg.Wait()
		close(msgChan)
	}()

	for v := range msgChan {
		msg += v
	}
	msg += fmt.Sprintf("所有任务已完成，一共%d个账号，成功登录%d个账号！\n",
		len(cfg.Accounts), countOpt.Load())
	Logger.Debug(msg)
	SendMessage(msg)
	Logger.Info("End terminal login")
}

func login(u models.Account, wg *sync.WaitGroup, msgChan chan<- string,
	countOpt *atomic.Int32) {
	defer wg.Done()
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
		withField.Error(err)
		msgChan <- fmt.Sprintf("%s: %s 登录失败(%s) %s \n\n", u.Addr,
			u.Username, utils.NowCST(), err)
		return
	}
	defer CloseSession(session)

	out, err := session.CombinedOutput("uname")
	if err != nil {
		withField.Errorf("登录失败: %s", err)
		msgChan <- fmt.Sprintf("%s: %s 登录失败(%s) %s \n\n", u.Addr,
			u.Username, utils.NowCST(), string(out))
		return
	}

	defer withField.Debug(string(out))
	withField.Info("SSH登录成功")
	msgChan <- fmt.Sprintf("%s: %s 登录成功(%s) %s \n", u.Addr,
		u.Username, utils.NowCST(), string(out))
	countOpt.Add(1)
}
