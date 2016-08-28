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
	"encoding/json"
	"errors"
	"flag"
	"sync"
	"time"
	// "fmt"

	"github.com/oikomi/FishChatServer/base"
	// "github.com/oikomi/FishChatServer/common"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/protocol"
	"github.com/oikomi/FishChatServer/storage/redis_store"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

type msgServerClientState struct {
	Session          *libnet.Session
	Alive            bool
	Valid            bool
	ClientSessionNum uint64
}
type PushServer struct {
	cfg             *PushServerConfig
	sessions        base.SessionMap
	channels        base.ChannelMap
	server          *libnet.Server
	voiceCache      *redis_store.VoiceCache
	sessionCache    *redis_store.SessionCache
	offlineMsgCache *redis_store.OfflineMsgCache
	p2pStatusCache  *redis_store.P2pStatusCache
	//	mysqlStore       *mysql_store.MysqlStore
	scanSessionMutex sync.Mutex
	readMutex        sync.Mutex // multi client session may ask for REDIS at the same time

	MsgServerClientMap   map[string]*msgServerClientState
	msgServerClientMutex sync.Mutex
}

func NewPushServer(cfg *PushServerConfig, rs *redis_store.RedisStore) *PushServer {
	return &PushServer{
		cfg:                cfg,
		MsgServerClientMap: make(map[string]*msgServerClientState), //已经连接的消息服务器表
		sessions:           make(base.SessionMap),
		channels:           make(base.ChannelMap),
		server:             new(libnet.Server),
		voiceCache:         redis_store.NewVoiceCache(rs),
		sessionCache:       redis_store.NewSessionCache(rs),
		offlineMsgCache:    redis_store.NewOfflineMsgCache(rs),
		p2pStatusCache:     redis_store.NewP2pStatusCache(rs),
		//		mysqlStore:      mysql_store.NewMysqlStore(cfg.Mysql.Addr, cfg.Mysql.Port, cfg.Mysql.User, cfg.Mysql.Password, cfg.Mysql.Database, cfg.Mysql.MaxOpenConn, cfg.Mysql.MaxOIdleConn),
	}
}

/*
   心跳检测消息服务器是否存活
*/
func (self *PushServer) scanDeadClient() {
	timer := time.NewTicker(180 * time.Second)
	for {
		select {
		case <-timer.C:
			log.Info("begin scanDeadClient")
			go func() {
				for ms, msgServerClient := range self.MsgServerClientMap {
					self.msgServerClientMutex.Lock()
					if msgServerClient.Alive == false || msgServerClient.Valid == false {
						//心跳没有收到回复，链接作废
						log.Warningf("CloseDeadClient [%s],Alive=%t,Valid=%t.", ms, msgServerClient.Alive, msgServerClient.Valid)
						msgServerClient.Session.Close()
					} else {
						//发送心跳，等待回复
						msgServerClient.Alive = false
						cmd := protocol.NewCmdSimple(protocol.SEND_PING_CMD)
						cmd.AddArg(self.cfg.UUID)
						err := msgServerClient.Session.Send(libnet.Json(cmd))
						if err != nil {
							msgServerClient.Session.Close()
						}
					}
					self.msgServerClientMutex.Unlock()
				}
			}()
		}
	}
}
func (self *PushServer) connectMsgServer(ms string) (*libnet.Session, error) {
	client, err := libnet.Dial("tcp", ms)

	return client, err
}

/*
   处理消息服务器发送过来的数据
*/

func (self *PushServer) handleMsgServerClient(msc *libnet.Session) error {
	err := msc.Process(func(msg *libnet.InBuffer) error {
		var c protocol.CmdSimple
		ms := msc.Conn().RemoteAddr().String()
		if self.MsgServerClientMap[ms] == nil {
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
			self.MsgServerClientMap[ms].Valid = true
		case protocol.PING_CMD_ACK:
			self.MsgServerClientMap[ms].Alive = true
		}
		return nil
	})
	return err
}

func (self *PushServer) subscribeChannels() error {
	self.msgServerClientMutex.Lock()
	defer self.msgServerClientMutex.Unlock()
	for _, ms := range self.cfg.MsgServerList {
		if self.MsgServerClientMap[ms] != nil {
			//已经创建过链接并且链接正常
			continue
		}
		msgServerClient, err := self.connectMsgServer(ms)
		if err != nil {
			log.Error(err.Error())
			go self.subscribeChannels()
			continue
		}
		//连接建立成功，开始发送通道订阅
		cmd := protocol.NewCmdSimple(protocol.SUBSCRIBE_CHANNEL_CMD)
		cmd.AddArg(protocol.SYSCTRL_PUSH_SERVER)
		cmd.AddArg(self.cfg.UUID)

		err = msgServerClient.Send(libnet.Json(cmd))
		if err != nil {
			log.Error(err.Error())
			go self.subscribeChannels()
			continue
		}
		//通道订阅发送成功
		self.MsgServerClientMap[ms] = new(msgServerClientState)
		self.MsgServerClientMap[ms].Alive = true
		self.MsgServerClientMap[ms].Valid = false
		self.MsgServerClientMap[ms].ClientSessionNum = 0
		self.MsgServerClientMap[ms].Session = msgServerClient

		//开始处理 消息服务器-> 接入服务器 的数据
		go func(ms string) {
			// go self.removeMsgServer(ms)
			err := self.handleMsgServerClient(msgServerClient)
			log.Infof("err=%s", err)
			if err != nil {
				// self.msgServerClientRWMutex.Lock()
				// defer self.msgServerClientRWMutex.Unlock()
				delete(self.MsgServerClientMap, ms)
				log.Info("delete ok")
			}
			go self.subscribeChannels()
		}(ms)
	}

	return nil
}
func (self *PushServer) createChannels() {
	log.Info("createChannels")
	for _, c := range base.ChannleList {
		channel := libnet.NewChannel(self.server.Protocol())
		self.channels[c] = base.NewChannelState(c, channel)
	}
}

func (self *PushServer) sendMonitorData() error {
	log.Info("sendMonitorData")
	resp := protocol.NewCmdMonitor()

	// resp.SessionNum = (uint64)(len(self.sessions))

	// log.Info(resp)

	mb := NewMonitorBeat("monitor", self.cfg.MonitorBeatTime, 40, 10)

	if self.channels[protocol.SYSCTRL_MONITOR] != nil {
		for {
			resp.SessionNum = (uint64)(len(self.sessions))

			//log.Info(resp)
			mb.Beat(self.channels[protocol.SYSCTRL_MONITOR].Channel, resp)
		}
		// _, err := self.channels[protocol.SYSCTRL_MONITOR].Channel.Broadcast(libnet.Json(resp))
		// if err != nil {
		// 	glog.Error(err.Error())
		// 	return err
		// }
	}

	return nil
}

func (self *PushServer) scanDeadSession() {
	log.Info("scanDeadSession")
	timer := time.NewTicker(self.cfg.ScanDeadSessionTimeout * time.Second)
	ttl := time.After(self.cfg.Expire * time.Second)
	for {
		select {
		case <-timer.C:
			log.Info("scanDeadSession timeout")
			go func() {
				for id, s := range self.sessions {
					self.scanSessionMutex.Lock()
					//defer self.scanSessionMutex.Unlock()
					if (s.State).(*base.SessionState).Alive == false {
						log.Info("delete" + id)
						self.procOffline(id)
					} else {
						s.State.(*base.SessionState).Alive = false
					}
					self.scanSessionMutex.Unlock()
				}
			}()
		case <-ttl:
			break
		}
	}
}

func (self *PushServer) procOnline(ID string) {
	// load all the topic list of this user
	sessionCacheData, err := self.sessionCache.Get(ID)
	if err != nil {
		log.Errorf("ID(%s) no session cache", ID)
		return
	}
	sessionCacheData.Alive = true
	self.sessionCache.Set(sessionCacheData)
}

func (self *PushServer) procOffline(ID string) {
	// load all the topic list of this user
	if self.sessions[ID] != nil {
		self.sessions[ID].Close()
		delete(self.sessions, ID)

		sessionCacheData, err := self.sessionCache.Get(ID)
		if err != nil {
			log.Errorf("ID(%s) no session cache", ID)
			return
		}
		sessionCacheData.Alive = false
		self.sessionCache.Set(sessionCacheData)
	}
}

func (self *PushServer) parseProtocol(cmd []byte, session *libnet.Session) error {
	var c protocol.CmdSimple
	// receive msg, that means client alive
	if session.State != nil {
		self.scanSessionMutex.Lock()
		session.State.(*base.SessionState).Alive = true
		self.scanSessionMutex.Unlock()
	}
	err := json.Unmarshal(cmd, &c)
	if err != nil {
		log.Error("error:", err)
		return err
	}
	pp := NewProtoProc(self)

	self.readMutex.Lock()
	defer self.readMutex.Unlock()

	log.Info(c)

	switch c.GetCmdName() {
	case protocol.SEND_PING_CMD: //接入服务器心跳请求
		err = pp.procPing(&c, session)
		if err != nil {
			log.Error("error:", err)
			return err
		}
	case protocol.SUBSCRIBE_CHANNEL_CMD: //订阅通道
		err = pp.procSubscribeChannel(&c, session)
		if err != nil {
			log.Error("error:", err)
			return err
		}

	case protocol.ACTION_TRANSFER_TO_DEVICE: //转发指令到设备
		err = pp.procTransferToDevice(&c, session)
		if err != nil {
			return err
		}

	}

	return err
}
