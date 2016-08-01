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
	// "errors"
	// "fmt"
	// "strconv"

	"github.com/oikomi/FishChatServer/connect_base"
	// "github.com/oikomi/FishChatServer/common"
	"github.com/oikomi/FishChatServer/connect_libnet"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/protocol"
	// "github.com/oikomi/FishChatServer/storage/mongo_store"
	// "github.com/oikomi/FishChatServer/storage/redis_store"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

type ProtoProc struct {
	connectServer *ConnectServer
}

func NewProtoProc(connectServer *ConnectServer) *ProtoProc {
	return &ProtoProc{
		connectServer: connectServer,
	}
}

// func (self *ProtoProc) procSubscribeChannel(cmd protocol.Cmd, session *connect_libnet.Session) {
// 	log.Info("procSubscribeChannel")
// 	channelName := cmd.GetArgs()[0]
// 	cUUID := cmd.GetArgs()[1]
// 	log.Info(channelName)
// 	if self.connectServer.channels[channelName] != nil {
// 		self.connectServer.channels[channelName].Channel.Join(session, nil)
// 		self.connectServer.channels[channelName].ClientIDlist = append(self.connectServer.channels[channelName].ClientIDlist, cUUID)
// 	} else {
// 		log.Warning(channelName + " is not exist")
// 	}

// 	log.Info(self.connectServer.channels)
// }

func (self *ProtoProc) procGetMinLoadMsgServer() string {
	var minload uint64
	var minloadserver string
	var msgServer string

	minload = 0xFFFFFFFFFFFFFFFF

	for str, msgServerClient := range self.connectServer.msgServerClientMap {
		if minload > msgServerClient.ClientSessionNum && msgServerClient.Valid == true {
			minload = msgServerClient.ClientSessionNum
			minloadserver = str
		}
	}
	msgServer = minloadserver
	return msgServer
}

func (self *ProtoProc) procCheckMsgServer(session *connect_libnet.Session) error {
	reSelect := false
	MsgServer:=session.State.(*connect_base.SessionState).MsgServer
	if (	session.State.(*connect_base.SessionState).MsgServer=="" || 
			self.connectServer.msgServerClientMap[session.State.(*connect_base.SessionState).MsgServer] == nil || 
			self.connectServer.msgServerClientMap[session.State.(*connect_base.SessionState).MsgServer].Valid != true){
        self.connectServer.msgServerClientEmptyMutex.RLock()
        defer self.connectServer.msgServerClientEmptyMutex.RUnlock()
		//重新选择消息处理服务器
		session.State.(*connect_base.SessionState).MsgServer=self.procGetMinLoadMsgServer()
		reSelect=true
		if (session.State.(*connect_base.SessionState).MsgServer==MsgServer){
			reSelect=false
		}
	}

	if (reSelect){
		//告知消息服务器
		cmd := protocol.NewCmdSimple(protocol.ACTION_SELECT_MSG_SERVER_CMD)
		cmd.Infos["IMEI"]=session.IMEI
		self.procTransferMsgServer(cmd, session)
	}
	return nil
}

func (self *ProtoProc) procTransferMsgServer(cmd protocol.Cmd,session *connect_libnet.Session) error {
	log.Info("procTransferMsgServer")
	log.Info(cmd)
	// self.connectServer.msgServerClientRWMutex.RLock()
	// defer self.connectServer.msgServerClientRWMutex.RUnlock()
	log.Info("procTransferMsgServer begin")

	msgServer:=session.State.(*connect_base.SessionState).MsgServer
	msgServerClient:=self.connectServer.msgServerClientMap[msgServer]
	for {
		msgServer=session.State.(*connect_base.SessionState).MsgServer
		msgServerClient=self.connectServer.msgServerClientMap[msgServer]
		if msgServerClient==nil {
			self.procCheckMsgServer(session)
		}else{
			break
		}
	}
	err := msgServerClient.Session.Send(libnet.Json(cmd))
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
