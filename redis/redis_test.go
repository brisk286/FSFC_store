package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"testing"
)

func init() {
	//连接redis
	Rdb := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379", Password: "", DB: 0})
	//健康检测
	_, err = Rdb.Ping().Result()
	if err != nil {
		log.Fatalln("redis状态错误: ", err)
	}
}

func Test_Kv(t *testing.T) {
	//新增k/v
	_ = Rdb.Set("hello", "world", 0).Err()

	//获取k/v
	result, _ := Rdb.Get("hello").Result()
	fmt.Println(result)

	//删除
	//_, _ = rdb.Del(ctx, "hello").Result()
}

func Test_List(t *testing.T) {
	//新增
	_ = Rdb.RPush("list", "message").Err()
	_ = Rdb.RPush("list", "message2").Err()

	//查询
	result, _ := Rdb.LLen("list").Result()
	fmt.Println(result)

	//更新
	_ = Rdb.LSet("list", 2, "message set").Err()

	//遍历
	lRange, _ := Rdb.LRange("list", 0, result).Result()
	for _, v := range lRange {
		log.Println(v)
	}

	//删除
	_, _ = Rdb.LRem("list", 3, "message2").Result()
}
