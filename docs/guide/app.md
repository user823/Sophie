# sophie app框架

Sophie 所有组件都基于pkg/app 框架进行构建。以gateway为例，该组件分为两个部分：
 - cmd/gateway/main.go 作为sophie-gateway的入口
 - internal/gateway 作为网关服务的主体部分

## app 配置方式
Sophie app 框架基于cobra，pflag，通过viper从命令行、配置文件、etcd配置中心中读取配置
每个组件基本都包含：
1. app.go 用于创建运行程序
2. options.go 用于该组件的配置选项
3. config.go 配置选项确定后，在完善各种参数后生成app实际运行的config
4. server.go 它定义了利用config.go 获取app对外提供服务的server

在pkg/app/options.go 中定义了4个接口，它们影响app的启动流程：
 - CliOptions: 实现该接口的option 表示该option会提供命令行flag
 - ValidatableOptions: 实现该接口的option 表示该option会在app启动前进行参数验证
 - CompletableOptions: 实现该接口的option 表示该option会在app启动前自行补全参数配置
 - PrintableOptions: 实现该接口的option 表示该option会在app启动前打印自己的参数

以gateway/options.go为例:
```
type Options struct {
	ServiceDiscover    *options.ServiceDiscoverOptions `json:"server_discover" mapstructure:"server_discover"`
	ServerRunOptions   *options.GenericRunOptions      `json:"generic" mapstructure:"generic"`
	InsecureServing    *options.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`
	SecureServing      *options.SecureServingOptions   `json:"secure" mapstructure:"secure"`
	RedisOptions       *options.RedisOptions           `json:"redis" mapstructure:"redis"`
	RPCClient          *options.RPCClientOptions       `json:"rpc_client" mapstructure:"rpc_client"`
	Log                *log.Options                    `json:"log" mapstructure:"log"`
	Jwt                *options.JwtOptions             `json:"jwt" mapstructure:"jwt"`
	AggregationOptions *aggregation.AnalyticsOptions   `json:"aggregation" mapstructure:"aggregation"`
	Availability       *options.AvailabilityOptions    `json:"availability" mapstructure:"availability"`
}
```
不管从那种来源获取配置，mapstructure 对应的属性会在app启动前绑定到对应字段上，因此写yml配置文件时要对应mapstructure

## app 启动流程
app 通过buildCommand完成自身构建，然后进入执行流程，整体来看：
[app 构建流程](../images/app.png)

构建app后，app运行时执行的runCommand 实际上对应每个组件的app.go 的run方法
run方法依次执行：创建server运行的配置、初始化log、创建server
[server 执行流程](../images/server.png)

