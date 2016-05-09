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
	"flag"

	"github.com/oikomi/FishChatServer/common"
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

type ProtoProc struct {
	gateway *Gateway
}

func NewProtoProc(gateway *Gateway) *ProtoProc {
	return &ProtoProc{
		gateway: gateway,
	}
}

// select a min load msg server
func (self *ProtoProc) procGetMinLoadMsgServer() string {
	var minload uint64
	var minloadserver string
	var msgServer string

	minload = 0xFFFFFFFFFFFFFFFF

	for str, msn := range self.gateway.msgServerNumMap {
		if minload > msn {
			minload = msn
			minloadserver = str
		}
	}
	msgServer = minloadserver
	return msgServer
}

func (self *ProtoProc) procLogin(cmd protocol.Cmd, session *libnet.Session) error {
	//log.Info("procLogin")
	var err error
	var uuid string
	var msgServer string

	ClientID := cmd.GetArgs()[0]
	ClientType := cmd.GetArgs()[1]
	ClientPwd := cmd.GetArgs()[2]

	// get the session cache
	sessionCacheData, err := self.gateway.sessionCache.Get(ClientID)
	if sessionCacheData != nil {
		log.Warningf("ID %s already login", ClientID)

		msgServer = sessionCacheData.MsgServerAddr
		uuid = sessionCacheData.ID
	} else {
		// choose msg server and allocate UUID
		msgServer = self.procGetMinLoadMsgServer()
		uuid = common.NewV4().String()
		// get the session store to check whether registered
		sessionStoreData, _ := self.gateway.mongoStore.GetSessionFromCid(ClientID)
		if sessionStoreData == nil {
			log.Warningf("ID %s not registered", ClientID)

			// for store data
			sessionStoreData = mongo_store.NewSessionStoreData(ClientID, ClientPwd, ClientType)
			log.Info(sessionStoreData)
			common.StoreData(self.gateway.mongoStore, sessionStoreData)
		}
		// for cache data, MsgServer MUST update local & remote addr.
		sessionCacheData = redis_store.NewSessionCacheData(sessionStoreData, session.Conn().RemoteAddr().String(), msgServer, uuid)
		log.Info(sessionCacheData)
		common.StoreData(self.gateway.sessionCache, sessionCacheData)
	}
	//
	resp := protocol.NewCmdSimple(protocol.RSP_LOGIN_CMD)
	resp.AddArg(protocol.RSP_SUCCESS)
	resp.AddArg(uuid)
	resp.AddArg(msgServer)

	log.Info("Resp | ", resp)

	if session != nil {
		err = session.Send(libnet.Json(resp))
		if err != nil {
			log.Error(err.Error())
		}
		session.Close()
		log.Info("client ", session.Conn().RemoteAddr().String(), " | close")
	}
	return nil
}
