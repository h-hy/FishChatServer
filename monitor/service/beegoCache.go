package service

import (
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/oikomi/FishChatServer/log"
)

var redisCache cache.Cache

func init() {
	var err error
	redisCache, err = cache.NewCache("redis", `{"key":"ApiServer","conn":"127.0.0.1:6379","dbNum":"0","password":""}`)
	if err != nil {
		log.Fatal(err)
	}
}

func GetBeegoCache() cache.Cache {
	return redisCache
}
