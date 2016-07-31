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
	// "encoding/json"
	"flag"
	"sync"
	"time"
	"fmt"
	"bytes"
	"errors"

	"github.com/oikomi/FishChatServer/connect_base"
	// "github.com/oikomi/FishChatServer/common"
	"github.com/oikomi/FishChatServer/connect_libnet"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/protocol"
	"github.com/oikomi/FishChatServer/storage/mongo_store"
	"github.com/oikomi/FishChatServer/storage/redis_store"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

type ConnectServer struct {
	cfg              *ConnectServerConfig
	sessions         connect_base.SessionMap
	channels         connect_base.ChannelMap
	topics           protocol.TopicMap
	server           *connect_libnet.Server
	sessionCache     *redis_store.SessionCache
	topicCache       *redis_store.TopicCache
	offlineMsgCache  *redis_store.OfflineMsgCache
	p2pStatusCache   *redis_store.P2pStatusCache
	mongoStore       *mongo_store.MongoStore
	scanSessionMutex sync.Mutex
	readMutex        sync.Mutex // multi client session may ask for REDIS at the same time
}

func NewMsgServer(cfg *ConnectServerConfig, rs *redis_store.RedisStore) *ConnectServer {
	return &ConnectServer{
		cfg:             cfg,
		sessions:        make(connect_base.SessionMap),
		channels:        make(connect_base.ChannelMap),
		topics:          make(protocol.TopicMap),
		server:          new(connect_libnet.Server),
		sessionCache:    redis_store.NewSessionCache(rs),
		topicCache:      redis_store.NewTopicCache(rs),
		offlineMsgCache: redis_store.NewOfflineMsgCache(rs),
		p2pStatusCache:  redis_store.NewP2pStatusCache(rs),
		mongoStore:      mongo_store.NewMongoStore(cfg.Mongo.Addr, cfg.Mongo.Port, cfg.Mongo.User, cfg.Mongo.Password),
	}
}

func (self *ConnectServer) createChannels() {
	log.Info("createChannels")
	for _, c := range connect_base.ChannleList {
		channel := connect_libnet.NewChannel(self.server.Protocol())
		self.channels[c] = connect_base.NewChannelState(c, channel)
	}
}

func (self *ConnectServer) sendMonitorData() error {
	// log.Info("sendMonitorData")
	// resp := protocol.NewCmdMonitor()

	// // resp.SessionNum = (uint64)(len(self.sessions))

	// // log.Info(resp)

	// mb := NewMonitorBeat("monitor", self.cfg.MonitorBeatTime, 40, 10)

	// if self.channels[protocol.SYSCTRL_MONITOR] != nil {
	// 	for {
	// 		resp.SessionNum = (uint64)(len(self.sessions))

	// 		//log.Info(resp)
	// 		mb.Beat(self.channels[protocol.SYSCTRL_MONITOR].Channel, resp)
	// 	}
	// 	// _, err := self.channels[protocol.SYSCTRL_MONITOR].Channel.Broadcast(connect_libnet.Json(resp))
	// 	// if err != nil {
	// 	// 	glog.Error(err.Error())
	// 	// 	return err
	// 	// }
	// }

	return nil
}

func (self *ConnectServer) scanDeadSession() {
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
					if (s.State).(*connect_base.SessionState).Alive == false {
						log.Info("delete" + id)
						self.procOffline(id)
					} else {
						s.State.(*connect_base.SessionState).Alive = false
					}
					self.scanSessionMutex.Unlock()
				}
			}()
		case <-ttl:
			break
		}
	}
}

func (self *ConnectServer) procOffline(ID string) {
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

var info_spilt []byte  = []byte{'|'}

var data_spilt []byte  = []byte{'#'}
var data_begin []byte  = []byte{'{','<'}
var data_end []byte  = []byte{'>','}'}
func (self *ConnectServer) parseProtocol(cmd []byte, session *connect_libnet.Session) error {
	var c protocol.CmdSimple
	fmt.Printf("parseProtocol\n\n")

	var index,nowindex,last_index int
	index = bytes.Index(cmd,data_begin)
	if (index==-1){
		return errors.New("cmd ERROR")
	}
	nowindex=index+1
	//数据头定位完毕
	infos := bytes.Split(cmd[0:index],info_spilt)
	if (len(infos) != 4 ){
		return errors.New("cmd ERROR")
	}
	c.IMEI=string(infos[3])
	//IMEI提取完毕
	last_index = bytes.Index(cmd,data_end)
	if (index==-1){
		return errors.New("cmd ERROR")
	}
	//数据尾定位完毕


	index = bytes.Index(cmd[nowindex:],data_spilt)
	if (index==-1){
		return errors.New("cmd ERROR")
	}
	c.CmdName = string(cmd[nowindex+1:nowindex+index])
	nowindex+=index+1
	//命令编码提取完毕
	for {
		index = bytes.Index(cmd[nowindex:],data_spilt)
		if (index==-1){
			if (nowindex != last_index){
				arg :=string(cmd[nowindex:last_index])
				c.Args = append(c.Args, arg)
			}
			break;
		}
		arg :=string(cmd[nowindex:nowindex+index])
		c.Args = append(c.Args, arg)
		nowindex+=index+1
	}
	//完整提取完毕

	// pp := NewProtoProc(self)

	self.readMutex.Lock()
	defer self.readMutex.Unlock()
	var err error
	log.Infof("[%s]->[%s]", session.Conn().RemoteAddr().String(), self.cfg.LocalIP)
	log.Info(c)
	
	
	// switch c.GetCmdName() {

	// }

	return err
}
