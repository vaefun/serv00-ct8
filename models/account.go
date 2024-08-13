package models

type Account struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Panel    string `mapstructure:"panel"`
	Addr     string `mapstructure:"addr"`
	IsCheck  bool   `mapstructure:"is_check"`
}
