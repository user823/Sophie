package doc

import (
	"crypto/tls"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/user823/Sophie/pkg/utils"
	"net"
	"net/http"
	"time"
)

// 配置elasticsearch 客户端
type ESConfig struct {
	Addrs         []string
	Username      string
	Password      string
	APIKey        string
	CloudID       string
	MaxIdle       int
	MaxRetryTimes int
	UseSSL        bool
	CA            string
	Timeout       time.Duration
}

func NewES(config any) (*elasticsearch.TypedClient, error) {
	esCfg, ok := config.(*ESConfig)
	if !ok {
		return nil, fmt.Errorf("config is not valid: %v", config)
	}

	retryBackoff := backoff.NewExponentialBackOff()
	return elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses:     esCfg.Addrs,
		Username:      esCfg.Username,
		Password:      esCfg.Password,
		APIKey:        esCfg.APIKey,
		CloudID:       esCfg.CloudID,
		RetryOnStatus: []int{502, 503, 504, 429},
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},
		MaxRetries: esCfg.MaxRetryTimes,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   esCfg.MaxIdle,
			ResponseHeaderTimeout: esCfg.Timeout,
			DialContext:           (&net.Dialer{Timeout: esCfg.Timeout}).DialContext,
			TLSHandshakeTimeout:   esCfg.Timeout,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: esCfg.UseSSL == false,
				RootCAs:            utils.LoadCAs(esCfg.CA),
			},
		},
	})
}
