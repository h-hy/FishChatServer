package models

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/oikomi/FishChatServer/log"
)

type User struct {
	Id        int `orm:"auto;pk;column(id)"`
	Username  string
	Password  string
	Telephone string
	Ticket    string
	Openid    string
	Devices   []*Device `orm:"rel(m2m);rel_through(github.com/oikomi/FishChatServer/monitor/models.UserDevice)"`
	isCache   bool      `orm:"-;default(false)"`
}

//func (u *User) TableName() string {
//	return "users"
//}

func UserRegist(username, telephone, password, openid string) (int, string, map[string]interface{}) {

	o := orm.NewOrm()
	user := User{Username: username}

	err := o.Read(&user, "Username")
	log.Info(err)
	if err == nil {
		return 20002, "用户名已经存在", map[string]interface{}{}
	} else if err == orm.ErrNoRows {
		newTicket := GetNewTicket()
		var newUser User
		newUser.Username = username
		newUser.Telephone = username
		newUser.Ticket = newTicket
		newUser.Password = password
		newUser.Openid = openid

		_, err := o.Insert(&newUser)
		if err == nil {
			newUser.CacheUser()
			return 0, "注册成功", map[string]interface{}{
				"username": username,
				"ticket":   newTicket,
			}
		}
		log.Info(err)
	}
	return 20006, "注册失败，请与管理员联系", map[string]interface{}{}
}
func (u *User) CacheUser() {
	log.Info("CacheUser")
	body, err := json.Marshal(u)
	log.Info(err)
	if err == nil {
		redisCache.Put("user_"+u.Username, body, 30*time.Minute)
	}
}

func (u *User) UpdateDevice() {
	o := orm.NewOrm()
	o.LoadRelated(u, "Devices")
	u.CacheUser()
}
func UserCleanOpenid(openid string) {
	o := orm.NewOrm()
	o.QueryTable("user").Filter("openid", openid).Update(orm.Params{
		"openid": "",
	})
}
func GetUser(username string) (User, error) {
	user := redisCache.Get("user_" + username)
	if user == nil {
		//缓存没有，从数据库读取
		o := orm.NewOrm()
		user := User{Username: username}

		err := o.Read(&user, "Username")
		log.Info(err)
		if err == nil {
			_, err = o.LoadRelated(&user, "Devices")
			if err != nil {
				log.Info(err)
				return user, err
			}
			user.isCache = false
			user.CacheUser()
		}
		return user, err
	} else {
		userString := GetString(user)
		var user User
		err := json.Unmarshal([]byte(userString), &user)
		user.isCache = true
		return user, err
	}
}
func (u *User) CheckTicket(ticket string) bool {
	return u.Ticket == ticket && u.Ticket != ""
}
func (u *User) CheckBind(IMEI string) bool {

	for _, device := range u.Devices {
		if device.IMEI == IMEI {
			return true
		}
	}
	if u.isCache == true {
		//重新获取一次试试
		//		redisCache.Delete("user_" + username)
		u.UpdateDevice()
		log.Info(IMEI)
		for _, device := range u.Devices {
			log.Info(device.IMEI)
			if device.IMEI == IMEI {
				return true
			}
		}
	}
	return false
}

//func UserCheckTicket(username, ticket string) bool {
//	log.Info(ticket)
//	if ticket == "" {
//		return false
//	}
//	user, err := GetUser(username)
//	log.Info(user)
//	log.Info(ticket)
//	if err == nil {
//		return user.Ticket == ticket && user.Ticket != ""
//	} else {
//		return false
//	}
//}

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

func init() {
	orm.RegisterModel(new(User))
}
