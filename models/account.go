package models

type Account struct {
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
	Panel    string `mapstructure:"panel" json:"panel"`
	Addr     string `mapstructure:"addr" json:"addr"`
	IsCheck  bool   `mapstructure:"is_check" json:"is_check"`
}
