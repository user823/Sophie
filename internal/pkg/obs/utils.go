package obs

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"github.com/user823/Sophie/pkg/utils"
	"github.com/user823/Sophie/pkg/utils/hash"
	"github.com/user823/Sophie/pkg/utils/strutil"
)

var (
	// 设置文件白名单
	AllowedPrefix = []string{".jpg", ".jpeg", ".png", ".gif", ".mp4", ".txt", ".bmp"}
	timeFormat    = "20060102150405"
)

func CheckFileName(fileName string) bool {
	prefix := filepath.Ext(fileName)
	return strutil.ContainsAny(prefix, AllowedPrefix...)
}

// 生成obs访问路径
// 命名规则:hash(userid + fileName) - "yy-MM-dd-ss" + prefix
func GetNewFileName(fileName string, userId int64) string {
	hasher := hash.NewHasher(hash.DefaultHashAlgorithm)
	prefix := filepath.Ext(fileName)

	createdTime := time.Now().Format(timeFormat)
	fileName = fmt.Sprintf("%d%s", userId, fileName)
	return hasher.HashKey(utils.S2b(fileName)) + "-" + createdTime + prefix
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
	bucketName := viper.GetString(driver + ".bucket")

	endpoint := viper.GetString(driver + ".endpoint")
	u := &url.URL{
		Scheme: schema,
		Host:   endpoint,
		Path:   path.Join(bucketName, objectName),
	}
	return u.String()
}
