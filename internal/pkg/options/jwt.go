package options

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"time"
)

// 对jwt选项进行配置
type JwtOptions struct {
	Realm      string        `json:"realm"       mapstructure:"realm"`
	Key        string        `json:"key"         mapstructure:"key"`
	Timeout    time.Duration `json:"timeout"     mapstructure:"timeout"`
	MaxRefresh time.Duration `json:"max-refresh" mapstructure:"max-refresh"`
}

func NewJwtOptions() *JwtOptions {
	return &JwtOptions{
		Realm:      "sophie jwt",
		Timeout:    1 * time.Hour,
		MaxRefresh: 1 * time.Hour,
	}
}

func (o *JwtOptions) Validate() error {
	if len(o.Key) < 6 || len(o.Key) > 32 {
		return fmt.Errorf("--sercret-key must between 6 and 32")
	}
	return nil
}

func (o *JwtOptions) AddFlags(fs *flag.FlagSet) {
	if fs == nil {
		return
	}
	fs.StringVar(&o.Realm, "jwt.realm", o.Realm, "Realm name to display to the user.")
	fs.StringVar(&o.Key, "jwt.key", o.Key, "Private key used to sign jwt token.")
	fs.DurationVar(&o.Timeout, "jwt.timeout", o.Timeout, "JWT token timeout.")

	fs.DurationVar(&o.MaxRefresh, "jwt.max-refresh", o.MaxRefresh, ""+
		"This field allows clients to refresh their token until MaxRefresh has passed.")
}
