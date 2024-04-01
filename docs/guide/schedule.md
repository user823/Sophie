# shedule

sophie 任务调度系统使用 "github.com/robfig/cron/v3"第三方库，它与quartz部分兼容
它使用标准的[cron 表达式](https://en.wikipedia.org/wiki/Cron)

前端新建定时任务信息（系统监控 -> 定时任务）
任务名称：自定义，如：定时查询任务状态
任务分组：根据字典sys_job_group配置
调用目标字符串：调用后台预定义的类名方法及其参数, 格式: 方法(参数1，参数2...)
执行表达式：标准cron表达式
状态：是否启动定时任务
备注：定时任务描述信息

调用目标： 系统启动时加载预定义的定时任务

注意：调用目标仅支持基本类型，不然不能正确解析

参考：
[cron](https://pkg.go.dev/github.com/robfig/cron/v3#section-readme)