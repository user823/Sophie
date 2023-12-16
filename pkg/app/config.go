package app

import (
	"fmt"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
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
		// 从配置文件中加载
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
		} else {
			viper.AddConfigPath(".")
			viper.AddConfigPath("../configs")

			if names := strings.Split(name, "-"); len(names) > 1 {
				viper.AddConfigPath(filepath.Join("/etc", names[0]))
				if homeDir := os.Getenv("HOME"); homeDir != "" {
					viper.AddConfigPath(filepath.Join(homeDir, names[0]))
				}
			}
			viper.SetConfigName(name)
		}

		// 从配置中心中加载
		viper.AddRemoteProvider("etcd3", viper.GetString("etcd3"), "/config")
		viper.AddRemoteProvider("etcd", viper.GetString("etcd"), "/config")
		viper.SetConfigName(name)
		viper.SetConfigType("json") // because there is no file extension in a stream of bytes, supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"

		if err := viper.ReadInConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to read configuration with configuration file or configuration center (%s): %v\n", cfgFile, err.Error())
			//os.Exit(1)
		}
	})
	return flag.Lookup(configName)
}

// 运行App基本默认配置
func SetDefaultConfig() {
	// 配置中心
	viper.SetDefault("etcd", "http://127.0.0.1:4001")
	viper.SetDefault("etcd3", "http://127.0.0.1:4001")
}
