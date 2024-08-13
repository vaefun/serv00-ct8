package models

// Config 定义了配置文件的结构
type Config struct {
	UUID             string    `mapstructure:"UUID"`
	Port             int       `mapstructure:"PORT"`
	LogLevel         string    `mapstructure:"log_level" json:"log_level"`
	PrivateKey       string    `mapstructure:"private_key" json:"private_key"`
	PushPlusToken    string    `mapstructure:"push_plus_token" json:"push_plus_token"`
	TelegramBotToken string    `mapstructure:"telegram_bot_token" json:"telegram_bot_token"`
	TelegramChatId   string    `mapstructure:"telegram_chat_id" json:"telegram_chat_id"`
	Accounts         []Account `mapstructure:"accounts" json:"accounts"`
}
