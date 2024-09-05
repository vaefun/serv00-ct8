package service

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/arlettebrook/serv00-ct8/configs"
	"github.com/arlettebrook/serv00-ct8/models"
	"github.com/arlettebrook/serv00-ct8/utils"

	"github.com/pkg/errors"

	"golang.org/x/crypto/ssh"
)

const keyName = "id_rsa"

// ConnSSH 根据私钥连接服务器
func ConnSSH(u models.Account) (c *ssh.Client, err error) {
	withField := Logger.WithField(u.Addr, u.Username)
	withField.Info("Start connection...")

	config := &ssh.ClientConfig{
		User: u.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(u.Password),
			ssh.KeyboardInteractive(func(name, instruction string,
				questions []string, echos []bool) (answers []string, err error) {
				answers = make([]string, len(questions))
				for i := range answers {
					answers[i] = u.Password
				}
				return answers, nil
			}),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if configs.Cfg.PrivateKey != "" {
		// 加载私钥
		key := []byte(configs.Cfg.PrivateKey)

		// 创建公钥签名器
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			withField.Errorf("unable to parse private key: %v", err)
			return nil, fmt.Errorf("unable to parse private key: %w", err)
		}
		config.Auth = append(config.Auth, ssh.PublicKeys(signer))
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", u.Addr),
		config)
	if err != nil {
		withField.Errorf("Failed to dial: %v", err)
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	withField.Info("SSH连接成功")

	return client, nil
}

// InitSSH 自动生成公钥、私钥，并上传公钥
func InitSSH() {
	Logger.Info("Start init ssh...")
	generateKey()
	saveKey()
	sendKey()
	Logger.Warn("公钥如果存在上传失败，请稍后在试，或者手动上传公钥，没有请忽略")
	Logger.Info("Init SSH end")
}

func generateKey() {
	if !utils.FileIsExists(keyName) {
		cmd := exec.Command("ssh-keygen", "-t", "rsa", "-f", keyName,
			"-C", "github@actions", "-N", "")
		output, err := cmd.CombinedOutput()
		if err != nil {
			Logger.Fatalf("CombinedOutput error: %s", err)
		}

		Logger.Debugf("Run ssh-keygen: %s", string(output))
	}
	Logger.Info("Generate key successfully")
}

func saveKey() {
	privateFile, err := os.ReadFile(keyName)
	if err != nil {
		Logger.Fatalf("Read key file %s error: %s", keyName, err)
	}
	pubFile, err := os.ReadFile(keyName + ".pub")
	if err != nil {
		Logger.Fatalf("Read key file %s error: %s", keyName, err)
	}

	msg := fmt.Sprintf(
		"以下内容为公钥：\n%s\n\n请将以下内容保存到PRIVATE_KEY环境变量中: \n%s",
		string(pubFile), string(privateFile))
	Logger.Warn("必须配置一个推送渠道，否则无法获取PRIVATE_KEY，已配置，请忽略")
	SendMessage(msg)
}

func sendKey() {
	var wg sync.WaitGroup
	msgChan := make(chan string, 3)
	msg := "所有账号公钥上传结果如下: \n"
	for _, u := range cfg.Accounts {
		wg.Add(1)
		go func(u models.Account) {
			defer wg.Done()
			execSendKey(u, msgChan)
		}(u)
	}

	go func() {
		wg.Wait()
		close(msgChan)
	}()

	for v := range msgChan {
		msg += v
	}
	msg += "如果失败的较少，建议手动上传失败主机的公钥或稍后重试\n"
	SendMessage(msg)
	Logger.Debug(msg)
}

func execSendKey(u models.Account, msgChan chan<- string) {
	withField := Logger.WithField(u.Addr, u.Username)
	withField.Info("Start upload...")
	cmd := exec.Command("bash", "-c",
		fmt.Sprintf("sshpass -p '%s' ssh-copy-id -o "+
			"StrictHostKeyChecking=no -i '%s' '%s@%s'",
			u.Password, keyName, u.Username, u.Addr))

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		withField.Errorf("公钥上传失败")
		withField.Debugf("Exec send key error: %s", err)
		msgChan <- fmt.Sprintf("%s: %s %s \n", u.Addr, u.Username,
			"公钥上传失败")
		return
	}

	defer withField.Debug(out.String())

	withField.Info("公钥上传成功")
	msgChan <- fmt.Sprintf("%s: %s %s \n", u.Addr, u.Username,
		"公钥上传成功")

}

func CloseClient(c *ssh.Client) {
	if err := c.Close(); err != nil {
		Logger.Warnf("Close client error: %s", err)
	}
}

func CloseSession(session *ssh.Session) {
	if err := session.Close(); err != nil {
		if errors.Is(err, io.EOF) {
			// todo：solve io.EOF
			return
		}
		Logger.Warnf("Close session error: %s", err)
	}
}

func CreateSession(c *ssh.Client) (s *ssh.Session, err error) {
	session, err := c.NewSession()
	if err != nil {
		Logger.Errorf("创建会话失败: %s", err)
		return nil, fmt.Errorf("创建会话失败: %s", err)
	}
	return session, nil
}
