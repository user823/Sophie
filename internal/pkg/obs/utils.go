package obs

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/user823/Sophie/pkg/utils/hash"
	"github.com/user823/Sophie/pkg/utils/strutil"
	"net/url"
	"path/filepath"
	"time"
)

var (
	// 设置文件白名单
	allowedPrefix = []string{".jpg", ".jpeg", ".png", ".gif", ".mp4", ".txt"}
	timeFormat    = "20060102150405"
)

func CheckFileName(fileName string) bool {
	prefix := filepath.Ext(fileName)
	return strutil.ContainsAny(prefix, allowedPrefix...)
}

// 生成obs访问路径
// 命名规则:hash(userid + fileName) - "yy-MM-dd-ss" + prefix
func GetNewFileName(fileName string, userId int64) string {
	hasher := hash.NewHasher(hash.DefaultHashAlgorithm)
	prefix := filepath.Ext(fileName)

	createdTime := time.Now().Format(timeFormat)
	fileName = fmt.Sprintf("%d%s", userId, fileName)
	return hasher.HashKey(fileName) + "-" + createdTime + prefix
}

// 拼接访问的url
func GetURL(objectName string) string {
	driver := viper.GetString("driver")
	if driver == "" {
		return ""
	}
	schema := "https"
	if !viper.GetBool(driver + ".use_ssl") {
		schema = "http"
	}

	endpoint := viper.GetString(driver + ".endpoint")
	u := &url.URL{
		Scheme: schema,
		Host:   endpoint,
		Path:   objectName,
	}
	return u.String()
}
