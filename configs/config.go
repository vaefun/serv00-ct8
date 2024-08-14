package configs

import (
	"encoding/json"
	"os"

	"github.com/arlettebrook/serv00-ct8/models"
	"github.com/spf13/viper"
)

const (
	cfgName         = "config.json"
	defaultLogLevel = "info"
)

var Cfg *models.Config

func init() {
	viper.SetConfigFile(cfgName)

	viper.SetDefault("log_level", defaultLogLevel)

	if err := viper.ReadInConfig(); err != nil {
		if !os.IsNotExist(err) {
			panic("Loading config error: " + err.Error())
		}
		// todo：配置文件config.json不存在，不会警告
	}

	assertBindEnvErr(viper.BindEnv("private_key"))
	assertBindEnvErr(viper.BindEnv("push_plus_token"))
	assertBindEnvErr(viper.BindEnv("telegram_bot_token"))
	assertBindEnvErr(viper.BindEnv("telegram_chat_id"))
	assertBindEnvErr(viper.BindEnv("accounts_json"))
	assertBindEnvErr(viper.BindEnv("log_level"))

	var cfg models.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic("Unmarshal config error: " + err.Error())
	}

	// 绑定环境变量获取的是string类型数据，需要手动序列化。
	accountsJson := viper.GetString("ACCOUNTS_JSON")
	if accountsJson != "" {
		err := json.Unmarshal([]byte(accountsJson), &cfg)
		if err != nil {
			panic("Unmarshal account_json error: " + err.Error())
		}
	}

	Cfg = &cfg
}

func assertBindEnvErr(err error) {
	if err != nil {
		panic("BindEnv error: " + err.Error())
	}
}
