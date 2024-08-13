package cmd

import (
	"github.com/arlettebrook/serv00-ct8/service"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "生成公钥、私钥。",
	Run: func(cmd *cobra.Command, args []string) {
		service.InitSSH()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
