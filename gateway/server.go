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

	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/protocol"
	"github.com/oikomi/FishChatServer/storage/mongo_store"
	"github.com/oikomi/FishChatServer/storage/redis_store"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

type Gateway struct {
	cfg                *GatewayConfig
	msgServerClientMap map[string]*libnet.Session
	// client to each msg server getting their payload
	msgServerNumMap map[string]uint64
	// payload of each msg server
	server       *libnet.Server
	sessionCache *redis_store.SessionCache
	mongoStore   *mongo_store.MongoStore
}

func NewGateway(cfg *GatewayConfig, rs *redis_store.RedisStore) *Gateway {
	return &Gateway{
		cfg:                cfg,
		msgServerClientMap: make(map[string]*libnet.Session),
		msgServerNumMap:    make(map[string]uint64),
		server:             new(libnet.Server),
		sessionCache:       redis_store.NewSessionCache(rs),
		mongoStore:         mongo_store.NewMongoStore(cfg.Mongo.Addr, cfg.Mongo.Port, cfg.Mongo.User, cfg.Mongo.Password),
	}
}

func (self *Gateway) parseProtocol(cmd []byte, session *libnet.Session) error {
	var c protocol.CmdSimple
	err := json.Unmarshal(cmd, &c)
	if err != nil {
		log.Error("error:", err)
		return err
	}

	pp := NewProtoProc(self)

	switch c.GetCmdName() {
	case protocol.REQ_LOGIN_CMD:
		err = pp.procLogin(&c, session)
		if err != nil {
			log.Error("error:", err)
			return err
		}
	}

	return err
}

func (self *Gateway) handleMsgServerClient(msc *libnet.Session) {
	msc.Process(func(msg *libnet.InBuffer) error {
		log.Info("msg_server", msc.Conn().RemoteAddr().String(), " say: ", string(msg.Data))
		var c protocol.CmdMonitor

		err := json.Unmarshal(msg.Data, &c)
		if err != nil {
			log.Error("error:", err)
			return err
		}
		self.msgServerNumMap[msc.Conn().RemoteAddr().String()] = c.SessionNum
		return nil
	})
}

func (self *Gateway) connectMsgServer(ms string) (*libnet.Session, error) {
	client, err := libnet.Dial("tcp", ms)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}

	return client, err
}

func (self *Gateway) subscribeChannels() error {
	log.Info("gateway start to subscribeChannels")
	for _, ms := range self.cfg.MsgServerList {
		msgServerClient, err := self.connectMsgServer(ms)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		cmd := protocol.NewCmdSimple(protocol.SUBSCRIBE_CHANNEL_CMD)
		cmd.AddArg(protocol.SYSCTRL_MONITOR)
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
