package cmd

import (
	"github.com/arlettebrook/serv00-ct8/service"
	"github.com/spf13/cobra"
)

var checkPM2cmd = &cobra.Command{
	Use:   "check-pm2",
	Short: "检查并恢复PM2进程。",
	Run: func(cmd *cobra.Command, args []string) {
		service.CheckPM2()
	},
}

func init() {
	rootCmd.AddCommand(checkPM2cmd)
}
