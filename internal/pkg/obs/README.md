# obs

使用了 "github.com/eleven26/goss/goss" 第三方库，支持阿里云、华为云、minio等对象存储服务
安装：`go get -u github.com/eleven26/goss`
参考：[goss](https://pkg.go.dev/github.com/eleven26/goss#section-readme)

配置方式：
在configs目录下的settings.yml 中配置各种obs，比如minio:
```yaml
minio:
  endpoint: 10.211.55.3:9000
  bucket: sophie
  access_key: sophie
  secret_key: 12345678
  user_ssl: false

# 选择驱动
driver: minio
```

