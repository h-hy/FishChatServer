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
	"errors"
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
	if (	session.State.(*connect_base.SessionState).MsgServer=="" || 
			self.connectServer.msgServerClientMap[session.State.(*connect_base.SessionState).MsgServer] == nil || 
			self.connectServer.msgServerClientMap[session.State.(*connect_base.SessionState).MsgServer].Valid != true){
		session.State.(*connect_base.SessionState).MsgServer=self.procGetMinLoadMsgServer()
	}

	if (	session.State.(*connect_base.SessionState).MsgServer=="" || 
			self.connectServer.msgServerClientMap[session.State.(*connect_base.SessionState).MsgServer] ==nil || 
			self.connectServer.msgServerClientMap[session.State.(*connect_base.SessionState).MsgServer].Valid != true){
		return errors.New("No MsgServer Valid.")
	}
	return nil
}

func (self *ProtoProc) procTransferMsgServer(cmd protocol.Cmd,session *connect_libnet.Session) error {
	log.Info("procTransferMsgServer")
	log.Info(cmd)
	err := self.connectServer.msgServerClientMap[session.State.(*connect_base.SessionState).MsgServer].Session.Send(libnet.Json(cmd))
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

// func (self *ProtoProc) procGoOffLine(cmd protocol.Cmd, session *connect_libnet.Session) error {
// 	var c protocol.CmdSimple
// 	c.Infos=make(map[string]string)
// 	c.Infos["ID"]=string(infos[0])
// 	c.Infos["Project"]=string(infos[1])
// 	c.Infos["Version"]=string(infos[2])
// 	c.Infos["IMEI"]=string(infos[3])
// 	return nil
// }