//
// Copyright 2014 Hong Miao (miaohong@miaohong.org). All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/oikomi/FishChatServer/base"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/storage/redis_store"
)

const VERSION string = "0.10"

func init() {
	flag.Set("alsologtostderr", "false")
	flag.Set("log_dir", "false")
}

func version() {
	fmt.Printf("push_server version %s Copyright (c) 2014 Harold Miao (miaohong@miaohong.org)  \n", VERSION)
}

//var InputConfFile = flag.String("conf_file", "push_server.json", "input conf file name")

func handleSession(ms *PushServer, session *libnet.Session) {
	session.Process(func(msg *libnet.InBuffer) error {
		err := ms.parseProtocol(msg.Data, session)
		if err != nil {
			log.Error(err.Error())
		}

		return nil
	})
}

func main() {
	version()
	flag.Parse()
	cfg := NewPushServerConfig("push_server.json")
	err := cfg.LoadConfig()
	if err != nil {
		log.Error(err.Error())
		return
	}

	rs := redis_store.NewRedisStore(&redis_store.RedisStoreOptions{
		Network:        "tcp",
		Address:        cfg.Redis.Addr + cfg.Redis.Port,
		ConnectTimeout: time.Duration(cfg.Redis.ConnectTimeout) * time.Millisecond,
		ReadTimeout:    time.Duration(cfg.Redis.ReadTimeout) * time.Millisecond,
		WriteTimeout:   time.Duration(cfg.Redis.WriteTimeout) * time.Millisecond,
		Database:       1,
		KeyPrefix:      base.COMM_PREFIX,
	})

	orm.Debug = true
	orm.RegisterDriver("mysql", orm.DRMySQL)

	maxIdle := 30
	maxConn := 30

	mysqlUsername := cfg.Mysql.User
	mysqlPassword := cfg.Mysql.Password
	mysqlDatabase := cfg.Mysql.Database

	orm.RegisterDataBase("default", "mysql", mysqlUsername+":"+mysqlPassword+"@/"+mysqlDatabase+"?charset=utf8", maxIdle, maxConn)

	ms := NewPushServer(cfg, rs)

	ms.server, err = libnet.Listen(cfg.TransportProtocols, cfg.Listen)
	if err != nil {
		panic(err)
	}
	log.Info("push_server running at  ", ms.server.Listener().Addr().String())

	ms.createChannels()

	ms.subscribeChannels()

	go ms.scanDeadSession()

	go ms.sendMonitorData()

	ms.server.Serve(func(session *libnet.Session) {
		log.Info("a new client ", session.Conn().RemoteAddr().String(), " | come in")
		go handleSession(ms, session)
	})
}
