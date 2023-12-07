# Sophie-log

日志说明：
Sophie-log 将日志信息分为业务日志和系统日志，它们都通过Sophie-log 日志聚合系统进行收集、分析、处理
log包基于zap进行的二次开发，主要用于对接Sophie-log系统的Log Facade.
特点：
1. 支持不同级别的日志输出：`Debug`、`Info`、`Warn`、`Error`、`Panic`、`Fatal`
2. 支持多端输出：标准输出、文件等
3. 兼容标准库log输出
4. 支持结构化输出和文本输出
5. 可自定义配置
6. 支持打印调用者信息、调用栈帧信息

默认不开启日志聚合模式，日志聚合使用：
1. 可以在命令行中通过 --aggregation 开启日志聚合
2. 可以调用SetAggregation API 来开启日志聚合

如果使用日志聚合，为了与Sophie-log子模块结合必须附带必要的环境信息，调用WithValues添加环境信息
