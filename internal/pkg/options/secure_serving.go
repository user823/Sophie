package options

import (
	"crypto/tls"
	"fmt"

	flag "github.com/spf13/pflag"
	"github.com/user823/Sophie/pkg/utils"
)

// 支持https api
// 默认开启https服务
type SecureServingOptions struct {
	BindAddress string `json:"bind_address" mapstructure:"bind_address"`
	BindPort    int    `json:"bind_port"    mapstructure:"bind_port"`
	// Required 表示BindPort不能为空
	Required   bool
	ServerCert CertKey `json:"tls" mapstructure:"tls"`
}

type CertKey struct {
	// 证书文件
	Cert string `json:"cert_file" mapstructure:"cert_file"`
	// 私钥文件
	Key string `json:"private_key_file" mapstructure:"private_key_file"`
}

func NewSecureServingOptions() *SecureServingOptions {
	return &SecureServingOptions{
		BindAddress: "0.0.0.0",
		BindPort:    8081,
		Required:    true,
		// 证书和私钥默认存放位置：
		ServerCert: CertKey{
			Cert: "/var/run/sophie/sophie-cert.pem",
			Key:  "/var/run/sophie/sophie-key.pem",
		},
	}
}

func (o *SecureServingOptions) GenerateTLSConfig() *tls.Config {
	return &tls.Config{
		Certificates:       []tls.Certificate{utils.LoadX509KeyPair(o.ServerCert.Cert, o.ServerCert.Key)},
		InsecureSkipVerify: false,
	}
}

func (o *SecureServingOptions) Validate() error {
	if o == nil {
		return nil
	}
	if !utils.IsValidIP(o.BindAddress) {
		return fmt.Errorf("Error secure bind address %s, please use ipv4 or ipv6 ", o.BindAddress)
	}
	if o.Required {
		if o.BindPort < 1 || o.BindPort > 65535 {
			return fmt.Errorf("Error secure bind port %v must be between 1 and 65535", o.BindPort)
		}

		if !utils.FileExists(o.ServerCert.Cert) || !utils.FileExists(o.ServerCert.Key) {
			return fmt.Errorf("Error CertFile or KeyFile not exists: %s %s", o.ServerCert.Cert, o.ServerCert.Key)
		}
	}
	return nil
}

func (o *SecureServingOptions) AddFlags(fs *flag.FlagSet) {
	if fs == nil {
		return
	}
	fs.StringVar(&o.BindAddress, "secure.bind_address", o.BindAddress, ""+
		"The IP address on which to listen for the tls port. The "+
		"associated interface(s) must be reachable by the rest of the engine, and by CLI/web "+
		"clients. If blank, all interfaces will be used (0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")
	desc := "The port on which to serve HTTPS with authentication and authorization."
	if o.Required {
		desc += " It cannot be switched off with 0."
	} else {
		desc += " If 0, don't serve HTTPS at all."
	}
	fs.IntVar(&o.BindPort, "secure.bind_port", o.BindPort, desc)
	fs.StringVar(&o.ServerCert.Cert, "secure.tls.cert_file", o.ServerCert.Cert, ""+
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated "+
		"after server cert).")

	fs.StringVar(&o.ServerCert.Key, "secure.tls.private_key_file",
		o.ServerCert.Key, ""+
			"File containing the default x509 private key matching --secure.tls.cert-key.cert-file.")
}
