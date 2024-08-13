package cmd

import (
	"github.com/arlettebrook/serv00-ct8/service"

	"github.com/spf13/cobra"
)

var Logger = service.Logger

var rootCmd = &cobra.Command{
	Use:   "serv00-ct8",
	Short: "同时支持serv00与ct8自动化批量保号、保活。",
	Long: `同时支持serv00与ct8自动化批量保号、保活。
	具体支持自动登录面板、自动SSH登录、PM2进程恢复，
	并支持将结果推送到Telegram Bot、PushPlus微信公众号。`,
	Run: func(cmd *cobra.Command, args []string) {
		Logger.Info("Main running...")
		// todo：节点运行
		Logger.Warn("待完成节点运行...")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		Logger.Fatalf("Root command execute error: %s", err)
	}
}
