package models

// Config 定义了配置文件的结构
type Config struct {
	UUID             string    `mapstructure:"UUID"`
	Port             int       `mapstructure:"PORT"`
	LogLevel         string    `mapstructure:"log_level"`
	PrivateKey       string    `mapstructure:"private_key"`
	PushPlusToken    string    `mapstructure:"push_plus_token"`
	TelegramBotToken string    `mapstructure:"telegram_bot_token"`
	TelegramChatId   string    `mapstructure:"telegram_chat_id"`
	Accounts         []Account `mapstructure:"accounts"`
}
