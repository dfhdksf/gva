package config

type CodeConfig struct {
	Redis struct {
		Addr string `mapstructure:"addr" json:"addr" yaml:"addr"`
	} `mapstructure:"redis" json:"redis" yaml:"redis"`
	TemplateId     string `mapstructure:"template-id" json:"template-id" yaml:"template-id"`             // 短信模板ID
	Expiration     int    `mapstructure:"expiration" json:"expiration" yaml:"expiration"`                // 验证码过期时间(秒)
	ResendInterval int    `mapstructure:"resend-interval" json:"resend-interval" yaml:"resend-interval"` // 重发间隔(秒)
}
