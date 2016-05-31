//
// Copyright 2014 Hong Miao. All Rights Reserved.
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
	"sync"

	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/protocol"
	"github.com/oikomi/FishChatServer/storage/mongo_store"
	"github.com/oikomi/FishChatServer/storage/redis_store"
)

type Router struct {
	cfg                *RouterConfig
	msgServerClientMap map[string]*libnet.Session
	sessionCache       *redis_store.SessionCache
	topicCache         *redis_store.TopicCache
	mongoStore         *mongo_store.MongoStore
	topicServerMap     map[string]string
	readMutex          sync.Mutex
}

func NewRouter(cfg *RouterConfig, rs *redis_store.RedisStore) *Router {
	return &Router{
		cfg:                cfg,
		msgServerClientMap: make(map[string]*libnet.Session),
		sessionCache:       redis_store.NewSessionCache(rs),
		topicCache:         redis_store.NewTopicCache(rs),
		mongoStore:         mongo_store.NewMongoStore(cfg.Mongo.Addr, cfg.Mongo.Port, cfg.Mongo.User, cfg.Mongo.Password),
		topicServerMap:     make(map[string]string),
	}
}

func (self *Router) connectMsgServer(ms string) (*libnet.Session, error) {
	client, err := libnet.Dial("tcp", ms)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}

	return client, err
}

func (self *Router) handleMsgServerClient(msc *libnet.Session) {
	msc.Process(func(msg *libnet.InBuffer) error {
		log.Info("msg_server", msc.Conn().RemoteAddr().String(), " say: ", string(msg.Data))
		var c protocol.CmdInternal
		pp := NewProtoProc(self)
		err := json.Unmarshal(msg.Data, &c)
		if err != nil {
			log.Error("error:", err)
			return err
		}
		switch c.GetCmdName() {
		case protocol.REQ_SEND_P2P_MSG_CMD:
			err := pp.procSendMsgP2P(&c, msc)
			if err != nil {
				log.Warning(err.Error())
			}

		case protocol.IND_ACK_P2P_STATUS_CMD:
			err := pp.procAckP2pStatus(&c, msc)
			if err != nil {
				log.Warning(err.Error())
			}

		case protocol.REQ_SEND_TOPIC_MSG_CMD:
			err := pp.procSendMsgTopic(&c, msc)
			if err != nil {
				log.Warning(err.Error())
			}

		}
		return nil
	})
}

func (self *Router) subscribeChannels() error {
	log.Info("route start to subscribeChannels")
	for _, ms := range self.cfg.MsgServerList {
		msgServerClient, err := self.connectMsgServer(ms)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		cmd := protocol.NewCmdSimple(protocol.SUBSCRIBE_CHANNEL_CMD)
		cmd.AddArg(protocol.SYSCTRL_SEND)
		cmd.AddArg(self.cfg.UUID)

		err = msgServerClient.Send(libnet.Json(cmd))
		if err != nil {
			log.Error(err.Error())
			return err
		}

		cmd = protocol.NewCmdSimple(protocol.SUBSCRIBE_CHANNEL_CMD)
		cmd.AddArg(protocol.SYSCTRL_TOPIC_SYNC)
		cmd.AddArg(self.cfg.UUID)

		err = msgServerClient.Send(libnet.Json(cmd))
		if err != nil {
			log.Error(err.Error())
			return err
		}

		self.msgServerClientMap[ms] = msgServerClient
	}

	for _, msc := range self.msgServerClientMap {
		go self.handleMsgServerClient(msc)
	}
	return nil
}
