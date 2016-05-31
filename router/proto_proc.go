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

	//"github.com/oikomi/FishChatServer/common"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/protocol"
	//"github.com/oikomi/FishChatServer/storage/mongo_store"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

type ProtoProc struct {
	Router *Router
}

func NewProtoProc(r *Router) *ProtoProc {
	return &ProtoProc{
		Router: r,
	}
}

/*
   MsgServer -> Router
       REQ_SEND_P2P_MSG_CMD
       arg0: Sent2ID       //接收方用户ID
       arg1: Msg           //消息内容
       arg2: FromID        //发送方用户ID
       arg3: uuid          //MsgServer分配的消息uuid
*/
func (self *ProtoProc) procSendMsgP2P(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("procSendMsgP2P")
	var err error
	send2ID := cmd.GetArgs()[0]
	send2Msg := cmd.GetArgs()[1]
	log.Info(send2Msg)

	self.Router.readMutex.Lock()
	defer self.Router.readMutex.Unlock()
	cacheSession, err := self.Router.sessionCache.Get(send2ID)
	if cacheSession == nil || cacheSession.Alive == false {
		storeSession, err := self.Router.mongoStore.GetSessionFromCid(send2ID)
		if err != nil {
			log.Warningf("ID %s not registered, msg dropped", send2ID)
			return err
		}
		log.Info(storeSession)
		log.Warningf("ID registered but offline: %s", send2ID)

		cmd.ChangeCmdName(protocol.ROUTE_SEND_P2P_MSG_CMD)
		//ms := self.Router.cfg.MsgServerList[0]
		ms := session.Conn().RemoteAddr().String()
		err = self.Router.msgServerClientMap[ms].Send(libnet.Json(cmd))
		if err != nil {
			log.Error("error:", err)
			return err
		}
	} else {
		log.Info(cacheSession.MsgServerAddr)

		cmd.ChangeCmdName(protocol.ROUTE_SEND_P2P_MSG_CMD)
		log.Info(cmd)
		err = self.Router.msgServerClientMap[cacheSession.MsgServerAddr].Send(libnet.Json(cmd))
		if err != nil {
			log.Error("error:", err)
			return err
		}
	}

	return nil
}

func (self *ProtoProc) procAckP2pStatus(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("procAckP2pStatus")
	var err error

	send2ID := cmd.GetArgs()[2]

	self.Router.readMutex.Lock()
	defer self.Router.readMutex.Unlock()
	cacheSession, err := self.Router.sessionCache.Get(send2ID)
	if cacheSession == nil || cacheSession.Alive == false {
		storeSession, err := self.Router.mongoStore.GetSessionFromCid(send2ID)
		if err != nil {
			log.Warningf("ID %s not registered, msg dropped", send2ID)
			return err
		}
		log.Info(storeSession)
		log.Warningf("ID registered but offline: %s", send2ID)

		cmd.ChangeCmdName(protocol.ROUTE_ACK_P2P_STATUS_CMD)
		ms := session.Conn().RemoteAddr().String()
		err = self.Router.msgServerClientMap[ms].Send(libnet.Json(cmd))
		if err != nil {
			log.Error("error:", err)
			return err
		}
	} else {
		log.Info(cacheSession.MsgServerAddr)

		cmd.ChangeCmdName(protocol.ROUTE_ACK_P2P_STATUS_CMD)
		log.Info(cmd)
		err = self.Router.msgServerClientMap[cacheSession.MsgServerAddr].Send(libnet.Json(cmd))
		if err != nil {
			log.Error("error:", err)
			return err
		}
	}

	return nil
}

/*
   Router -> MsgServer
       ROUTE_SEND_TOPIC_MSG_CMD
       arg0: ClientID      //发送方用户ID
       arg1: ClientType    //发送方终端类型，是client还是device
       arg2: Msg           //消息内容
       arg3: TopicName     //群组名, device无须提供
*/
func (self *ProtoProc) procSendMsgTopic(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("procSendMsgTopic")
	var err error
	ms := session.Conn().RemoteAddr().String()

	//send2Msg := cmd.GetArgs()[0]
	topicName := cmd.GetArgs()[1]
	//fromID := cmd.GetArgs()[2]
	//fromType := cmd.GetArgs()[3]

	// check whether the topic exist
	topicCacheData, err := self.Router.topicCache.Get(topicName)
	if topicCacheData == nil {
		log.Warningf("TOPIC %s not exist: %s", topicName, err.Error())
		return err
	}

	cmd.ChangeCmdName(protocol.ROUTE_SEND_TOPIC_MSG_CMD)
	log.Info(protocol.ROUTE_SEND_TOPIC_MSG_CMD)
	log.Info(cmd)

	for ip, num := range topicCacheData.AliveMemberNumMap {
		if num > 0 {
			log.Warningf("topic %s has %d member(s) in ip %s", topicName, num, ip)
			if ip != ms {
				// not in this server, routing it
				err = self.Router.msgServerClientMap[ip].Send(libnet.Json(cmd))
				if err != nil {
					log.Error("error:", err)
					return err
				}
			}
		}
	}

	return nil
}
