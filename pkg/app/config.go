package app

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"github.com/user823/Sophie/pkg/log"
	"strings"
)

const (
	configName = "config"
)

var (
	cfgFile string
)

func init() {
	flag.StringVarP(&cfgFile, configName, "c", cfgFile, "Read configuration from specified `FILE`, "+
		"support JSON, TOML, YAML, HCL, or Java properties formats.")
}

func addConfigFlag(name string) *flag.Flag {
	// 读取环境变量
	viper.AutomaticEnv()
	viper.SetEnvPrefix(strings.Replace(strings.ToUpper(name), "-", "_", -1))
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	cobra.OnInitialize(func() {
		log.Debug("正在读取配置文件...")
		// 从配置文件中加载
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
		} else {
			viper.AddConfigPath(".")
			viper.AddConfigPath("./configs")
			viper.AddConfigPath("../configs")
			viper.AddConfigPath("../../configs")
			viper.SetConfigName(name)
		}
		viper.SetConfigType("yml")
		if err := viper.ReadInConfig(); err != nil {
			log.Warnf("failed to read configuration with configuration file (%s): %s", cfgFile, err.Error())
		}

		// 从配置文件中读取配置中心访问端点
		configCenter := viper.GetString("config_center.name")
		endpoint := viper.GetString("config_center.endpoint")
		if configCenter != "" && endpoint != "" {
			if err := viper.AddRemoteProvider(configCenter, endpoint, "/sophie-config/"+name+".json"); err != nil {
				log.Warnf("failed to add viper remoting config: %s", err.Error())
			} else {
				if err := viper.ReadRemoteConfig(); err != nil {
					log.Warnf("failed to read configuration from remoting: %s", err.Error())
					return
				}

				// 读取远程配置成功
				viper.Set("enable_remote", true)
			}
		}
	})
	return flag.Lookup(configName)
}
