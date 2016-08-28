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

package service

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/oikomi/FishChatServer/base"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/protocol"
	"github.com/oikomi/FishChatServer/storage/redis_store"
)

type Monitor struct {
	cfg                  *MonitorConfig
	SessionCache         *redis_store.SessionCache
	readMutex            sync.Mutex
	PushServerClientMap   map[string]*pushServerClientState
	pushServerClientMutex sync.Mutex
}

type pushServerClientState struct {
	Session          *libnet.Session
	Alive            bool
	Valid            bool
	ClientSessionNum uint64
}

func NewMonitor(cfg *MonitorConfig) *Monitor {
	return &Monitor{
		cfg:                cfg,
		PushServerClientMap: make(map[string]*pushServerClientState), //已经连接的消息服务器表
		SessionCache: redis_store.NewSessionCache(redis_store.NewRedisStore(&redis_store.RedisStoreOptions{
			Network:        "tcp",
			Address:        cfg.Redis.Addr + cfg.Redis.Port,
			ConnectTimeout: time.Duration(cfg.Redis.ConnectTimeout) * time.Millisecond,
			ReadTimeout:    time.Duration(cfg.Redis.ReadTimeout) * time.Millisecond,
			WriteTimeout:   time.Duration(cfg.Redis.WriteTimeout) * time.Millisecond,
			Database:       1,
			KeyPrefix:      base.COMM_PREFIX,
		})),
	}
}

/*
   心跳检测消息服务器是否存活
*/
func (self *Monitor) scanDeadClient() {
	timer := time.NewTicker(180 * time.Second)
	for {
		select {
		case <-timer.C:
			log.Info("begin scanDeadClient")
			go func() {
				for ms, pushServerClient := range self.PushServerClientMap {
					self.pushServerClientMutex.Lock()
					if pushServerClient.Alive == false || pushServerClient.Valid == false {
						//心跳没有收到回复，链接作废
						log.Warningf("CloseDeadClient [%s],Alive=%t,Valid=%t.", ms, pushServerClient.Alive, pushServerClient.Valid)
						pushServerClient.Session.Close()
					} else {
						//发送心跳，等待回复
						pushServerClient.Alive = false
						cmd := protocol.NewCmdSimple(protocol.SEND_PING_CMD)
						cmd.AddArg(self.cfg.UUID)
						err := pushServerClient.Session.Send(libnet.Json(cmd))
						if err != nil {
							pushServerClient.Session.Close()
						}
					}
					self.pushServerClientMutex.Unlock()
				}
			}()
		}
	}
}
func (self *Monitor) connectPushServer(ms string) (*libnet.Session, error) {
	client, err := libnet.Dial("tcp", ms)

	return client, err
}

/*
   处理消息服务器发送过来的数据
*/

func (self *Monitor) handlePushServerClient(msc *libnet.Session) error {
	err := msc.Process(func(msg *libnet.InBuffer) error {
		var c protocol.CmdSimple
		ms := msc.Conn().RemoteAddr().String()
		if self.PushServerClientMap[ms] == nil {
			log.Error(ms + " not exist")
			return errors.New(ms + " not exist")
		}
		err := json.Unmarshal(msg.Data, &c)
		if err != nil {
			log.Error("error:", err)
			return err
		}

		log.Infof("c.GetCmdName()=%s\n\n", c.GetCmdName())
		switch c.GetCmdName() {
		case protocol.SUBSCRIBE_CHANNEL_CMD_ACK:
			self.PushServerClientMap[ms].Valid = true
		case protocol.PING_CMD_ACK:
			self.PushServerClientMap[ms].Alive = true
		}
		return nil
	})
	return err
}

func (self *Monitor) subscribeChannels() error {
	self.pushServerClientMutex.Lock()
	defer self.pushServerClientMutex.Unlock()
	for _, ms := range self.cfg.PushServerList {
		if self.PushServerClientMap[ms] != nil {
			//已经创建过链接并且链接正常
			continue
		}
		pushServerClient, err := self.connectPushServer(ms)
		if err != nil {
			log.Error(err.Error())
			go self.subscribeChannels()
			continue
		}
		//连接建立成功，开始发送通道订阅
		cmd := protocol.NewCmdSimple(protocol.SUBSCRIBE_CHANNEL_CMD)
		cmd.AddArg(protocol.SYSCTRL_API_SERVER)
		cmd.AddArg(self.cfg.UUID)

		err = pushServerClient.Send(libnet.Json(cmd))
		if err != nil {
			log.Error(err.Error())
			go self.subscribeChannels()
			continue
		}
		//通道订阅发送成功
		self.PushServerClientMap[ms] = new(pushServerClientState)
		self.PushServerClientMap[ms].Alive = true
		self.PushServerClientMap[ms].Valid = false
		self.PushServerClientMap[ms].ClientSessionNum = 0
		self.PushServerClientMap[ms].Session = pushServerClient

		//开始处理 消息服务器-> 接入服务器 的数据
		go func(ms string) {
			// go self.removePushServer(ms)
			err := self.handlePushServerClient(pushServerClient)
			log.Infof("err=%s", err)
			if err != nil {
				// self.pushServerClientRWMutex.Lock()
				// defer self.pushServerClientRWMutex.Unlock()
				delete(self.PushServerClientMap, ms)
				log.Info("delete ok")
			}
			go self.subscribeChannels()
		}(ms)
	}

	return nil
}
