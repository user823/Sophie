# Sophie-errors

errors 可用于构建客户端的返回消息
主要特点：
1. 兼容标准库的errors包
2. 定制了error内容的格式化输出
3. 为errors添加了栈帧信息输出功能
4. 通过errors包注册http响应码及消息类型