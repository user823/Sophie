package options

import (
	"fmt"

	flag "github.com/spf13/pflag"
	"github.com/user823/Sophie/pkg/utils"
)

type InsecureServingOptions struct {
	BindAddress string `json:"bind_address" mapstructure:"bind_address"`
	BindPort    int    `json:"bind_port"    mapstructure:"bind_port"`
}

func NewInsecureServingOptions() *InsecureServingOptions {
	return &InsecureServingOptions{
		BindAddress: "0.0.0.0",
		BindPort:    8080,
	}
}

func (o *InsecureServingOptions) Validate() error {
	if o == nil {
		return nil
	}
	if !utils.IsValidIP(o.BindAddress) {
		return fmt.Errorf("Error insecure bind address %s, please use ipv4 or ipv6 ", o.BindAddress)
	}
	if o.BindPort < 1 || o.BindPort > 65535 {
		return fmt.Errorf("Error insecure bind port %d must be between 1 and 65535", o.BindPort)
	}
	return nil
}

func (o *InsecureServingOptions) AddFlags(fs *flag.FlagSet) {
	if fs == nil {
		return
	}
	fs.StringVar(&o.BindAddress, "insecure.bind_address", o.BindAddress, ""+
		"The IP address on which to serve the --insecure.bind_port "+
		"(set to 0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")
	fs.IntVar(&o.BindPort, "insecure.bind_port", o.BindPort, ""+
		"The port on which to serve unsecured, unauthenticated access. It is assumed "+
		"that firewall rules are set up such that this port is not reachable from outside of "+
		"the deployed machine This is performed by nginx in the default setup. Set to zero to disable.")
}
