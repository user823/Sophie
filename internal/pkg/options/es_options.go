package options

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/user823/Sophie/pkg/db/doc"
	"time"
)

type ESOptions struct {
	Addrs         []string      `json:"addrs" mapstructure:"addrs"`
	Username      string        `json:"username" mapstructure:"username"`
	Password      string        `json:"password" mapstructure:"password"`
	APIKey        string        `json:"api_key" mapstructure:"api_key"`
	CloudId       string        `json:"cloud_id" mapstructure:"cloud_id"`
	MaxIdle       int           `json:"max_idle" mapstructure:"max_idle"`
	MaxRetryTimes int           `json:"max_retry_times" mapstructure:"max_retry_times"`
	UseSSL        bool          `json:"use_ssl" mapstructure:"use_ssl"`
	CA            string        `json:"ca" mapstructure:"ca"`
	Timeout       time.Duration `json:"timeout" mapstructure:"timeout"`
}

func NewESOptions() *ESOptions {
	return &ESOptions{
		Addrs:         []string{"https://127.0.0.1:9200"},
		Username:      "sophie",
		Password:      "123456",
		APIKey:        "",
		CloudId:       "",
		MaxIdle:       10,
		UseSSL:        false,
		CA:            "",
		Timeout:       3 * time.Second,
		MaxRetryTimes: 3,
	}
}

func (o *ESOptions) Validate() error {
	if o.UseSSL == true && o.CA == "" {
		return fmt.Errorf("can not find valid elasticsearch ca file")
	}
	if o.CloudId != "" && len(o.Addrs) != 0 {
		return fmt.Errorf("both Addresses and CloudID are set")
	}
	return nil
}

func (o *ESOptions) AddFlags(fs *flag.FlagSet) {
	fs.StringSliceVar(&o.Addrs, "es.addrs", o.Addrs, ""+
		"A set of elasticsearch endpoints(eg:http://localhost:9200)")
	fs.StringVar(&o.Username, "es.username", o.Username, ""+
		"Username for access to elasticsearch")
	fs.StringVar(&o.Password, "es.password", o.Password, ""+
		"Password for access to elasticsearch")
	fs.StringVar(&o.APIKey, "es.apiKey", o.APIKey, ""+
		"Elasticsearch apikey for authorization(when set, username/password will be ignored)")
	fs.StringVar(&o.CloudId, "es.cloudId", o.CloudId, ""+
		"Endpoint for the Elastic Service (https://elastic.co/cloud)")
	fs.IntVar(&o.MaxIdle, "es.maxIdle", o.MaxIdle, ""+
		"Max idle connections per-host")
	fs.BoolVar(&o.UseSSL, "es.useSSL", o.UseSSL, ""+
		"Use SSL/TLS to build elasticsearch connection")
	fs.DurationVar(&o.Timeout, "es.timeout", o.Timeout, ""+
		"Timeout for connection、request et al.")
	fs.IntVar(&o.MaxRetryTimes, "es.max_retry_times", o.MaxRetryTimes, ""+
		"Elasticsearch max try times")
}

func (o *ESOptions) BuildESConfig() *doc.ESConfig {
	return &doc.ESConfig{
		Addrs:         o.Addrs,
		Username:      o.Username,
		Password:      o.Password,
		APIKey:        o.APIKey,
		CloudID:       o.CloudId,
		MaxIdle:       o.MaxIdle,
		MaxRetryTimes: o.MaxRetryTimes,
		UseSSL:        o.UseSSL,
		CA:            o.CA,
		Timeout:       o.Timeout,
	}
}
