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
	// "time"

	// "github.com/oikomi/FishChatServer/connect_base"
	"github.com/oikomi/FishChatServer/protocol"
	"github.com/oikomi/FishChatServer/connect_libnet"
	"github.com/oikomi/FishChatServer/log"
	// "github.com/oikomi/FishChatServer/storage/redis_store"
)

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
const char* build_time(void) {
	static const char* psz_build_time = "["__DATE__ " " __TIME__ "]";
	return psz_build_time;
}
*/
// import "C"

// var (
// 	buildTime = C.GoString(C.build_time())
// )

// func BuildTime() string {
// 	return buildTime
// }

const VERSION string = "0.10"

func init() {
	flag.Set("alsologtostderr", "false")
	flag.Set("log_dir", "false")
}

func version() {
	fmt.Printf("connect_server version %s Copyright (c) 2014 Harold Miao (miaohong@miaohong.org)  \n", VERSION)
}

var InputConfFile = flag.String("conf_file", "connect_server.json", "input conf file name")

func handleSession(ms *ConnectServer, session *connect_libnet.Session) {
	fmt.Printf("handleSession")
	session.Process(func(msg *connect_libnet.InBuffer) error {
		fmt.Printf("callback parseProtocol")
		err := ms.parseProtocol(msg.Data, session)
		if err != nil {
			log.Error(err.Error())
		}
		return nil
	})
}

func main() {
	version()
	// fmt.Printf("built on %s\n", BuildTime())
	flag.Parse()
	cfg := NewConnectServerConfig(*InputConfFile)
	err := cfg.LoadConfig()
	if err != nil {
		log.Error(err.Error())
		return
	}

	// rs := redis_store.NewRedisStore(&redis_store.RedisStoreOptions{
	// 	Network:        "tcp",
	// 	Address:        cfg.Redis.Addr + cfg.Redis.Port,
	// 	ConnectTimeout: time.Duration(cfg.Redis.ConnectTimeout) * time.Millisecond,
	// 	ReadTimeout:    time.Duration(cfg.Redis.ReadTimeout) * time.Millisecond,
	// 	WriteTimeout:   time.Duration(cfg.Redis.WriteTimeout) * time.Millisecond,
	// 	Database:       1,
	// 	KeyPrefix:      connect_base.COMM_PREFIX,
	// })

	// ms := NewMsgServer(cfg, rs)
	ms := NewMsgServer(cfg)
	ms.msgServerClientEmptyMutex.Lock()
	ms.server, err = connect_libnet.Listen(cfg.TransportProtocols, cfg.Listen)  //监听服务端口等待客户端连接
	if err != nil {
		panic(err)
	}
	log.Info("connect_server running at  ", ms.server.Listener().Addr().String())

	ms.subscribeChannels() //连接到消息服务器
	
	go ms.scanDeadSession()	//清理无用session
	
	go ms.scanDeadClient()	//清理无用消息服务器

	// go ms.sendMonitorData()

	ms.server.Serve(func(session *connect_libnet.Session) {
		log.Info("a new client ", session.Conn().RemoteAddr().String(), " | come in")
		session.AddCloseCallback(ms, func() {
			//客户端下线，通知消息服务器
			log.Info("AddCloseCallback callback")
			if (session.IMEI != ""){
				// ms.scanSessionMutex.Lock()
				cmd := protocol.NewCmdSimple(protocol.ACTION_GO_OFFLINE_CMD)
				cmd.Infos["IMEI"]=session.IMEI
				pp := NewProtoProc(ms)
				err = pp.procCheckMsgServer(session)
				if err != nil{
					// ms.scanSessionMutex.Unlock()
					return 
				}
				err = pp.procTransferMsgServer(cmd, session)
				// ms.scanSessionMutex.Unlock()
			}
		})
		go handleSession(ms, session)
	})
}


