package main

import (
	"fmt"

	"github.com/Yan0613/TikTok/middleware/rabbitmq"

	"github.com/Yan0613/TikTok/config"
	"github.com/Yan0613/TikTok/dao"
	"github.com/Yan0613/TikTok/log/logger"
	"github.com/Yan0613/TikTok/middleware/redis"
	"github.com/Yan0613/TikTok/utils"
)

func main() {
	defer logger.Sync()
	r := NewRouter()
	initDeps()

	myConfig := config.HTTPServer()
	err := r.Run(fmt.Sprintf(":%d", myConfig.Port)) // listen and serve on listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		return
	}
}

func initDeps() {
	//初始化数据库连接
	dao.Init()
	rabbitmq.Init()
	rabbitmq.InitCommentMQ()
	rabbitmq.InitRelationMQ()
	rabbitmq.InitLikeRabbitMQ()
	redis.InitRedis()
	utils.InitWordsFilter()
}
