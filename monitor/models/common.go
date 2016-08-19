package models

import (
	"fmt"

	"github.com/astaxie/beego/cache"
	"github.com/oikomi/FishChatServer/monitor/service"
)

var redisCache cache.Cache

func init() {
	redisCache = service.GetBeegoCache()
}

func GetString(v interface{}) string {
	switch result := v.(type) {
	case string:
		return result
	case []byte:
		return string(result)
	default:
		if v != nil {
			return fmt.Sprintf("%v", result)
		}
	}
	return ""
}
