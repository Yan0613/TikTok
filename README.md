# TikTok
基于Gin，GORM，Redis, MySQL, RabbitMQ开发的简易版抖音后端
## 项目结构
```
TikTok
├── /config/ 配置文件包
├── /controller/ 控制器包
├── /dao/ 数据库访问
├── /middleware/ 中间件
│   ├── ffmpeg/ 视频截图
│   ├── jwt/ 鉴权
│   ├── rabbitmq/ 消息队列
│   ├── redis/ 缓存
├── /service/ 服务层
├── /utils/ 工具
├── .gitignore
├── /go.mod/
├── LICENSE
├── main.go
├── README.md
└── router.go
```
## 项目依赖
```go
require (
	github.com/BurntSushi/toml v1.3.2
	github.com/aliyun/aliyun-oss-go-sdk v2.2.8+incompatible
	github.com/aliyun/credentials-go v1.3.1
	github.com/brianvoe/gofakeit/v6 v6.23.1
	github.com/gin-contrib/pprof v1.4.0
	github.com/gin-gonic/gin v1.9.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/importcjj/sensitive v0.0.0-20200106142752-42d1c505be7b
	github.com/streadway/amqp v1.1.0
	go.uber.org/zap v1.25.0
	golang.org/x/crypto v0.9.0
	gorm.io/driver/mysql v1.5.1
	gorm.io/gorm v1.25.2
)
```
项目运行：

```go
go run main.go router.go
```
## 项目配置
配置文件：
```
[HTTPServer]
IP = "" // 服务器 IP
Port =  // 服务器端口号
[Database]
IP = "" // 数据库 IP
Port =  // 数据库端口号
Account = "" // 数据库用户名
Password = "" // 数据库密码
DatabaseName = "" // 数据库名
Protocol = ""
Charset = ""
ParseTime = 
TimeZone = ""
[OSS]
CredentialType = ""
CredentialRoleName = ""
Endpoint = { Internal = "oss-cn-beijing-internal.aliyuncs.com", External = "oss-cn-beijing.aliyuncs.com" }
BucketName = ""
[Redis]
RedisHost = "" // Redis IP
RedisPort =  // Redis 端口号
RedisPassword = "" // Redis 密码
[Rabbitmq]
RabbitmqHost = "" // 消息队列 IP
RabbitmqPort =  // 消息队列端口号
RabbitmqUsername = "" // 消息队列用户名
RabbitmqPassword = "" // 消息队列密码
```
