package controllers

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	//	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/monitor/service"
	"github.com/oikomi/FishChatServer/storage/redis_store"
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

func restReturn(errcode int, errmsg string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"errcode": errcode,
		"errmsg":  errmsg,
		"data":    data,
	}
}

func getCache(m *service.Monitor, IMEI string) (*redis_store.SessionCacheData, error) {

	sessionCacheData, err := m.SessionCache.Get(IMEI)
	return sessionCacheData, err
}

const (
	KC_RAND_KIND_NUM   = 0 // 纯数字
	KC_RAND_KIND_LOWER = 1 // 小写字母
	KC_RAND_KIND_UPPER = 2 // 大写字母
	KC_RAND_KIND_ALL   = 3 // 数字、大小写字母
)

func GetNewTicket() string {
	return string(Krand(32, KC_RAND_KIND_LOWER))
}
func Krand(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}
