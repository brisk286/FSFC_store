package redis

import (
	"github.com/go-redis/redis"
	"log"
)

var (
	Client *redis.Client
	err    error
)

func init() {
	//连接redis
	//Client = redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379", Password: "root", DB: 0})
	Client = redis.NewClient(&redis.Options{Addr: "124.70.57.7:6379", Password: "root", DB: 0})

	//健康检测
	_, err = Client.Ping().Result()
	if err != nil {
		log.Fatalln("redis状态错误: ", err)
	}
}
