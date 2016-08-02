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
	"flag"
	"sync"
	"time"
	// "fmt"

	"github.com/oikomi/FishChatServer/base"
	// "github.com/oikomi/FishChatServer/common"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/protocol"
	"github.com/oikomi/FishChatServer/storage/mysql_store"
	"github.com/oikomi/FishChatServer/storage/redis_store"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

type MsgServer struct {
	cfg              *MsgServerConfig
	sessions         base.SessionMap
	channels         base.ChannelMap
	topics           protocol.TopicMap
	server           *libnet.Server
	sessionCache     *redis_store.SessionCache
	topicCache       *redis_store.TopicCache
	offlineMsgCache  *redis_store.OfflineMsgCache
	p2pStatusCache   *redis_store.P2pStatusCache
	mysqlStore       *mysql_store.MysqlStore
	scanSessionMutex sync.Mutex
	readMutex        sync.Mutex // multi client session may ask for REDIS at the same time
}

func NewMsgServer(cfg *MsgServerConfig, rs *redis_store.RedisStore) *MsgServer {
	return &MsgServer{
		cfg:             cfg,
		sessions:        make(base.SessionMap),
		channels:        make(base.ChannelMap),
		topics:          make(protocol.TopicMap),
		server:          new(libnet.Server),
		sessionCache:    redis_store.NewSessionCache(rs),
		topicCache:      redis_store.NewTopicCache(rs),
		offlineMsgCache: redis_store.NewOfflineMsgCache(rs),
		p2pStatusCache:  redis_store.NewP2pStatusCache(rs),
		mysqlStore:      mysql_store.NewMysqlStore(cfg.Mysql.Addr, cfg.Mysql.Port, cfg.Mysql.User, cfg.Mysql.Password, cfg.Mysql.Database, cfg.Mysql.MaxOpenConn, cfg.Mysql.MaxOIdleConn),
	}
}

func (self *MsgServer) createChannels() {
	log.Info("createChannels")
	for _, c := range base.ChannleList {
		channel := libnet.NewChannel(self.server.Protocol())
		self.channels[c] = base.NewChannelState(c, channel)
	}
}

func (self *MsgServer) sendMonitorData() error {
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

func (self *MsgServer) scanDeadSession() {
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

func (self *MsgServer) procOnline(ID string) {
	// load all the topic list of this user
	sessionCacheData, err := self.sessionCache.Get(ID)
	if err != nil {
		log.Errorf("ID(%s) no session cache", ID)
		return
	}
	sessionCacheData.Alive = true
	self.sessionCache.Set(sessionCacheData)
}

func (self *MsgServer) procOffline(ID string) {
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
		for _, topicName := range sessionCacheData.TopicList {
			topicCacheData, _ := self.topicCache.Get(topicName)
			if topicCacheData != nil {
				// update AliveMemberNumMap[server]
				if v, ok := topicCacheData.AliveMemberNumMap[self.cfg.LocalIP]; ok {
					if v > 0 {
						topicCacheData.AliveMemberNumMap[self.cfg.LocalIP] = v - 1
					} else {
						topicCacheData.AliveMemberNumMap[self.cfg.LocalIP] = 0
					}
					self.topicCache.Set(topicCacheData)
				}
			}
		}
	}
}

func (self *MsgServer) parseProtocol(cmd []byte, session *libnet.Session) error {
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
		case protocol.SEND_PING_CMD:		//接入服务器心跳请求
			err = pp.procPing(&c, session)
			if err != nil {
				log.Error("error:", err)
				return err
			}
		case protocol.SUBSCRIBE_CHANNEL_CMD: 	//订阅通道
			err = pp.procSubscribeChannel(&c, session)
			if err != nil {
				log.Error("error:", err)
				return err
			}

		case protocol.ACTION_GO_OFFLINE_CMD:	//终端下线通知
			err = pp.procGoOffLine(&c, session)
			if err != nil {
				log.Error("error:", err)
				return err
			}

		case protocol.ACTION_SELECT_MSG_SERVER_CMD: //终端切换消息服务器通知
			err = pp.procSelectMsgServer(&c, session)
			if err != nil {
				log.Error("error:", err)
				return err
			}

		case "U"+protocol.DEIVCE_HEARTBEAT_CMD:			//心跳包（上行）
			err = pp.prochHeartbeat(&c, session)
			if err != nil {
				return err
			}
		case "U"+protocol.DEIVCE_TIME_SYNC_CMD:			//时间同步（上行）
			err = pp.prochTimeSync(&c, session)
			if err != nil {
				return err
			}
		case "U"+protocol.DEIVCE_LOCATON_CMD:			//定位数据（上行）
			err = pp.prochLocation(&c, session)
			if err != nil {
				return err
			}
		case "U"+protocol.DEIVCE_LINK_DESC_CMD:			//连接用途请求通知（上行）
			err = pp.prochLinkDesc(&c, session)
			if err != nil {
				return err
			}
		case "U"+protocol.DEIVCE_VOICE_READED_CMD:			//连接用途请求通知（上行）
			err = pp.prochVoiceReaded(&c, session)
			if err != nil {
				return err
			}
		case "U"+protocol.DEIVCE_LOW_POWER_CMD:			//连接用途请求通知（上行）
			err = pp.prochLowPower(&c, session)
			if err != nil {
				return err
			}
		case "U"+protocol.DEIVCE_SOS_CMD:			//连接用途请求通知（上行）
			err = pp.prochSOS(&c, session)
			if err != nil {
				return err
			}

		// case protocol.REQ_LOGOUT_CMD:
		// 	err = pp.procLogout(&c, session)
		// 	if err != nil {
		// 		log.Error("error:", err)
		// 		return err
		// 	}

		// case protocol.REQ_SEND_P2P_MSG_CMD:
		// 	err = pp.procSendMessageP2P(&c, session)
		// 	if err != nil {
		// 		log.Error("error:", err)
		// 		return err
		// 	}
		// case protocol.ROUTE_SEND_P2P_MSG_CMD:
		// 	err = pp.procRouteMessageP2P(&c, session)
		// 	if err != nil {
		// 		log.Error("error:", err)
		// 		return err
		// 	}

		// // p2p ack
		// case protocol.IND_ACK_P2P_STATUS_CMD:
		// 	err = pp.procP2pAck(&c, session)
		// 	if err != nil {
		// 		log.Error("error:", err)
		// 		return err
		// 	}
		// // p2p ack
		// case protocol.ROUTE_ACK_P2P_STATUS_CMD:
		// 	err = pp.procP2pAck(&c, session)
		// 	if err != nil {
		// 		log.Error("error:", err)
		// 		return err
		// 	}

		// case protocol.REQ_SEND_TOPIC_MSG_CMD:
		// 	err = pp.procSendTopicMsg(&c, session)
		// 	if err != nil {
		// 		log.Error("error:", err)
		// 		return err
		// 	}
		// case protocol.ROUTE_SEND_TOPIC_MSG_CMD:
		// 	err = pp.procRouteTopicMsg(&c, session)
		// 	if err != nil {
		// 		log.Error("error:", err)
		// 		return err
		// 	}

		// case protocol.REQ_CREATE_TOPIC_CMD:
		// 	err = pp.procCreateTopic(&c, session)
		// 	if err != nil {
		// 		log.Error("error:", err)
		// 		return err
		// 	}

		// case protocol.REQ_ADD_2_TOPIC_CMD:
		// 	err = pp.procAdd2Topic(&c, session)
		// 	if err != nil {
		// 		log.Error("error:", err)
		// 		return err
		// 	}

		// case protocol.REQ_KICK_TOPIC_CMD:
		// 	err = pp.procKickTopic(&c, session)
		// 	if err != nil {
		// 		log.Error("error:", err)
		// 		return err
		// 	}

		// case protocol.REQ_JOIN_TOPIC_CMD:
		// 	err = pp.procJoinTopic(&c, session)
		// 	if err != nil {
		// 		log.Error("error:", err)
		// 		return err
		// 	}

		// case protocol.REQ_QUIT_TOPIC_CMD:
		// 	err = pp.procQuitTopic(&c, session)
		// 	if err != nil {
		// 		log.Error("error:", err)
		// 		return err
		// 	}

		// case protocol.REQ_GET_TOPIC_LIST_CMD:
		// 	err = pp.procGetTopicList(&c, session)
		// 	if err != nil {
		// 		log.Error("error:", err)
		// 		return err
		// 	}

		// case protocol.REQ_GET_TOPIC_MEMBER_CMD:
		// 	err = pp.procGetTopicMember(&c, session)
		// 	if err != nil {
		// 		log.Error("error:", err)
		// 		return err
		// 	}
	}

	return err
}
