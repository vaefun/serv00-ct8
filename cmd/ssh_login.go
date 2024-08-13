package cmd

import (
	"github.com/arlettebrook/serv00-ct8/service"
	"github.com/spf13/cobra"
)

var sshLoginCmd = &cobra.Command{
	Use:   "ssh-login",
	Short: "SSH登录。",
	Run: func(cmd *cobra.Command, args []string) {
		service.TerminalLogin()
	},
}

func init() {
	rootCmd.AddCommand(sshLoginCmd)
}
