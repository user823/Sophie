package obs

import (
	"fmt"
	"github.com/eleven26/goss/core"
	"github.com/eleven26/goss/goss"
	"github.com/spf13/viper"
	"github.com/user823/Sophie/pkg/log"
	"io"
	"strings"
)

var (
	storage core.Storage
)

func init() {
	// 如果没有加载，则手动加载一次
	if !viper.IsSet("driver") {
		viper.SetConfigFile("../../../configs/goss.yml")
		if err := viper.ReadInConfig(); err != nil {
			log.Errorf("viper 读取obs配置失败: %s", err.Error())
			return
		}
	}
	gs, err := goss.NewWithViper(viper.GetViper())
	if err != nil {
		log.Errorf("OBS 服务加载失败，请检查配置: %s", err.Error())
		return
	}
	storage = gs.Storage
}

// 直接从io流中读取文件信息并上传，返回访问url
func Upload(objectName string, r io.Reader) (string, error) {
	// 首先检查文件名后缀
	if !CheckFileName(objectName) {
		return "", fmt.Errorf("文件类型不正确，仅允许上传这些类型: %s", strings.Join(allowedPrefix, ", "))
	}

	url := GetURL(objectName)
	if err := storage.Put(objectName, r); err != nil {
		return "", err
	}
	return url, nil
}
