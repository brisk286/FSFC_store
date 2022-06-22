package redis

import (
	"github.com/go-redis/redis"
	"log"
)

var (
	Rdb *redis.Client
	err error
)

func init() {
	//连接redis
	Rdb = redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379", Password: "", DB: 0})
	//健康检测
	_, err = Rdb.Ping().Result()
	if err != nil {
		log.Fatalln("redis状态错误: ", err)
	}
}
