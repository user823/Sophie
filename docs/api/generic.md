# 通用说明

Sophie 系统 API 严格遵循 REST 标准进行设计，使用 JWT Token 进行 API 认证，请求编码格式均为：`UTF-8`格式。

## 1. 公共参数

每个参数都属于不同的类型，根据参数位置不同，参数有如下类型：

- 请求头参数 (Header)：例如 `Content-Type: application/json`。
- 路径参数 (Path)：例如 `/user/{id}` 中的 id 参数就位于 path 中。
- 查询参数 (Quexf)：例如 `users?username=colin&username=james&value=`。
- 请求体参数 (Body)：例如 `{"metadata":{"name":"secretdemo"},"expires":0,"description":"admin secret"}`。

这些公共参数，是每个接口都需要传入的，在每个接口文档中，不再一一说明。

IAM API 接口公共参数如下：

| 参数名称      | 位置   | 类型   | 必选 | 描述                                                 |
| ------------- | ------ | ------ | ---- | ---------------------------------------------------- |
| Content-Type  | Header | String | 是   | 取值必须是 application/json 或者 multipart/form-data |
| Authorization | Header | String | 是   | JWT Token，值以 `Bearer` 开头                        |

## 2. 返回结果

所有 API 的处理结果包含成功和失败两种情况，它们返回的内容有所差异，但一律以 json 格式传输：

- 成功时，返回结果包含以下内容：
  1. X-Request-Id：位于 HTTP 返回请求头中，调用的请求 ID，用来唯一表示一次请求。
  2. HTTP 状态码：HTTP 状态码，成功的请求状态码永远是 200。
  3. 接口请求的数据：位于 HTTP 返回 Body 中，API 请求需要的返回数据，JSON 格式。
- 失败时，返回结果包含以下内容：
  1. X-Request-Id：位于 HTTP 返回请求头中，调用的请求 ID，用来唯一表示一次请求。
  2. HTTP 状态码：HTTP 状态码，不同的错误类型返回 HTTP 状态码不同，可能状态码为 400、401、403、404、500.
  3. 返回的错误信息格式：`{"code": xxx, "message": "xxx"}`其中 code 表示错误码，message 表示错误描述信息。

## 3. 返回参数类型

由于参数是 JSON 格式，因此支持的参数类型为 string、number、array、boolean、null、object。
为了支持 JSON 与 golang 结构体的转换，object 类型直接使用结构体名代替，number 定义为更精确的类型：`Int、Uint、Int8、Uint8、Int16、Uint16、Int32、Uint32、Int64、Uint64、Float、Float64`

## 4. 认证

本系统支持 Basic 认证方式和 Bearer 认证方式（使用 JWT token）

## 5. 请求方法

本 API 接口文档中请求方法格式为：`HTTP方法 请求路径`，例如请求方法为：`GET /v1/users`, 请求地址为：`wuliserver.com`，请求协议为：`HTTP`，则实际的请求格式为：`curl -XGET http://wuliserver.com/v1/users`

## 6. 错误码

Sophie 系统同时返回 2 类错误码：HTTP 状态码和业务错误码。
HTTP 状态码分为 3 类：

- 200：代表成功相应
- 4xx：响应失败，说明客户端发生错误
- 500：响应失败，说明服务端发生错误

**HTTP 状态码说明**

| 状态码 | 说明                                       |
| ------ | ------------------------------------------ |
| 200    | 成功响应                                   |
| 400    | 客户端发生错误，比如参数不合法、格式错误等 |
| 401    | 认证失败                                   |
| 403    | 授权失败（比如菜单管理权限）               |
| 404    | 页面或者资源不存在                         |
| 500    | 响应失败，说明服务端发生了错误             |

**业务码 说明**
参考：error_code.md

## 7. 其他说明

1. v1 标识了当前 api 版本，v2,v3... 类推
