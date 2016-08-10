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
	"strconv"
	"time"
	// "bytes"

	"github.com/oikomi/FishChatServer/base"
	// "github.com/oikomi/FishChatServer/common"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/protocol"
	"github.com/oikomi/FishChatServer/storage/redis_store"
    "github.com/oikomi/FishChatServer/provider"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

type ProtoProc struct {
	msgServer *MsgServer
}

func NewProtoProc(msgServer *MsgServer) *ProtoProc {
	return &ProtoProc{
		msgServer: msgServer,
	}
}

func (self *ProtoProc) procSubscribeChannel(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("procSubscribeChannel")
	channelName := cmd.GetArgs()[0]
	cUUID := cmd.GetArgs()[1]
	log.Info(channelName)
	if self.msgServer.channels[channelName] != nil {
		if (channelName==protocol.SYSCTRL_CONNECT_SERVER){
			self.msgServer.sessions[cUUID] = session
			self.msgServer.sessions[cUUID].State = base.NewSessionState(true, cUUID, "connect_server")
			resp := protocol.NewCmdSimple(protocol.SUBSCRIBE_CHANNEL_CMD_ACK)
			_ = session.Send(libnet.Json(resp))
		}
		self.msgServer.channels[channelName].Channel.Join(session, nil)
		self.msgServer.channels[channelName].ClientIDlist = append(self.msgServer.channels[channelName].ClientIDlist, cUUID)
	} else {
		log.Warning(channelName + " is not exist")
	}
	return nil
	// log.Info(self.msgServer.channels)
}

func (self *ProtoProc) procPing(cmd protocol.Cmd, session *libnet.Session) error {
	//log.Info("procPing")
	//cid := session.State.(*base.SessionState).ClientID
		//self.msgServer.sessions[cid].State.(*base.SessionState).Alive = true
	if session.State != nil {
        // fmt.Printf("session.State != nil")
		self.msgServer.scanSessionMutex.Lock()
		defer self.msgServer.scanSessionMutex.Unlock()
		resp := protocol.NewCmdSimple(protocol.PING_CMD_ACK)
		// log.Info(resp)
        // fmt.Printf("resp=",resp)
		err := session.Send(libnet.Json(resp))
		if err != nil {
			log.Error(err.Error())
			return err
		}
		session.State.(*base.SessionState).Alive = true
	}
	return nil
}

func (self *ProtoProc) procGoOffLine(cmd protocol.Cmd, session *libnet.Session) error {
	IMEI := cmd.GetInfos()["IMEI"]

	sessionCacheData, err := self.msgServer.sessionCache.Get(IMEI)
	if sessionCacheData != nil {
		sessionCacheData.MsgServerAddr=""
		if (sessionCacheData.Alive != false){
			self.msgServer.mysqlStore.DeviceUpdate(IMEI,"alive=false")
		}
		sessionCacheData.Alive=false
		sessionCacheData.MaxAge=600 * time.Second
		self.msgServer.sessionCache.Set(sessionCacheData)
	}else{
		self.msgServer.mysqlStore.DeviceUpdate(IMEI,"alive=false")
	}


	return err
}

func (self *ProtoProc) procSelectMsgServer(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("procSelectMsgServer")
	var msgServer string
	IMEI := cmd.GetInfos()["IMEI"]


	// resp := protocol.NewCmdSimple(protocol.RSP_LOGIN_CMD)
	// for cache data
	msgServer = self.msgServer.cfg.LocalIP

	session.State.(*base.SessionState).Devices[IMEI]=msgServer

	// get the session store to check whether registered
	deviceStoreData, _ := self.msgServer.mysqlStore.GetDeviceFromIMEI(IMEI)
	if deviceStoreData == nil {
		log.Warningf("IMEI %s not registered", IMEI)
		// for store data
		self.msgServer.mysqlStore.DeviceStore(IMEI,"")
		// log.Info(deviceStoreData)
	}
	// for cache data, MsgServer MUST update local & remote addr.
	sessionCacheData := redis_store.NewSessionCacheData(session.Conn().RemoteAddr().String(), msgServer, deviceStoreData)
	log.Info(sessionCacheData)
	self.msgServer.sessionCache.Set(sessionCacheData)
	//处理完成
	return nil
}
func (self *ProtoProc) closeSession(IMEI string, session *libnet.Session) error {
	resp := protocol.NewCmdSimple(protocol.ACTION_DO_CLOSE_SESSION_CMD)
	resp.Infos["IMEI"]=IMEI
	log.Info("Resp | ", resp)

	if session != nil {
		err := session.Send(libnet.Json(resp))
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}
func (self *ProtoProc) checkCache(IMEI string, session *libnet.Session) (*redis_store.SessionCacheData,error) {
	sessionCacheData, err := self.msgServer.sessionCache.Get(IMEI)
	if sessionCacheData == nil {
		log.Warningf("no cache IMEI : %s, err: %s", IMEI, err.Error())
		deviceStoreData, _ := self.msgServer.mysqlStore.GetDeviceFromIMEI(IMEI)
		if deviceStoreData == nil {
			// not registered
			self.closeSession(IMEI,session)
			log.Warningf("no store IMEI : %s, err: %s", IMEI, err.Error())
			return nil,err
		}else{
			sessionCacheData := redis_store.NewSessionCacheData(session.Conn().RemoteAddr().String(), self.msgServer.cfg.LocalIP, deviceStoreData)
			return sessionCacheData,nil
		}
	}
	return sessionCacheData,nil
}
func (self *ProtoProc) prochHeartbeat(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("prochHeartbeat")
	resp := protocol.NewCmdSimple("C"+cmd.GetCmdName()[1:])
	IMEI := cmd.GetInfos()["IMEI"]
	if len(cmd.GetArgs()) == 2 {
		energy,err := strconv.Atoi(cmd.GetArgs()[1])
		if err != nil{
			resp.AddArg("2")
		}else{
			resp.AddArg("1")
			self.checkCache(IMEI,session)
			sessionCacheData, err := self.msgServer.sessionCache.Get(IMEI)
			if (err!=nil){
				return err
			}
			if (sessionCacheData.Energy != energy || sessionCacheData.Alive != true){
				sessionCacheData.Energy = energy
				if (sessionCacheData.Alive != true){
					self.msgServer.mysqlStore.DeviceUpdate(IMEI,"alive=true")
				}
				sessionCacheData.Alive = true
				self.msgServer.sessionCache.Set(sessionCacheData)
			}
		}
	}else{
		resp.AddArg("2")
	}

	resp.Infos["IMEI"]=IMEI

	if session != nil {
		err := session.Send(libnet.Json(resp))
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}

func (self *ProtoProc) prochTimeSync(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("prochTimeSync")
	resp := protocol.NewCmdSimple("C"+cmd.GetCmdName()[1:])
	IMEI := cmd.GetInfos()["IMEI"]
	resp.AddArg("1")
	resp.AddArg(time.Now().Format("2006-01-02 15:04:05"))

	resp.Infos["IMEI"]=IMEI

	if session != nil {
		err := session.Send(libnet.Json(resp))
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}

func (self *ProtoProc) prochLocation(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("prochLocation")
	resp := protocol.NewCmdSimple("C"+cmd.GetCmdName()[1:])
	IMEI := cmd.GetInfos()["IMEI"]
	if len(cmd.GetArgs()) == 1 {
		var locationInfo provider.Location
		locationInfo.Parse(cmd.GetArgs()[0])
		log.Info("locationInfo=",locationInfo)
		log.Info("locationInfo.LocationData=",locationInfo.LocationData)
		resp.AddArg("1")
	}else{
		resp.AddArg("2")
	}

	resp.Infos["IMEI"]=IMEI

	if session != nil {
		err := session.Send(libnet.Json(resp))
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}

func (self *ProtoProc) prochLinkDesc(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("prochLinkDesc")
	resp := protocol.NewCmdSimple("C"+cmd.GetCmdName()[1:])
	IMEI := cmd.GetInfos()["IMEI"]
	if len(cmd.GetArgs()) == 2 {
		commId := cmd.GetArgs()[0]
		resp.AddArg(commId)
		resp.AddArg("1")
	}else if len(cmd.GetArgs()) == 1 {
		commId := cmd.GetArgs()[0]
		resp.AddArg(commId)
		resp.AddArg("2")
	}else{
		resp.AddArg("2")
	}

	resp.Infos["IMEI"]=IMEI

	if session != nil {
		err := session.Send(libnet.Json(resp))
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}

func (self *ProtoProc) prochVoiceReaded(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("prochVoiceReaded")
	resp := protocol.NewCmdSimple("C"+cmd.GetCmdName()[1:])
	IMEI := cmd.GetInfos()["IMEI"]
	if len(cmd.GetArgs()) == 1 {
		id := cmd.GetArgs()[0]
		resp.AddArg(id)
		resp.AddArg("1")
	}else{
		resp.AddArg("2")
	}

	resp.Infos["IMEI"]=IMEI

	if session != nil {
		err := session.Send(libnet.Json(resp))
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}

func (self *ProtoProc) prochLowPower(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("prochLowPower")
	resp := protocol.NewCmdSimple("C"+cmd.GetCmdName()[1:])
	IMEI := cmd.GetInfos()["IMEI"]
	resp.AddArg("1")
	resp.Infos["IMEI"]=IMEI

	if session != nil {
		err := session.Send(libnet.Json(resp))
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}

func (self *ProtoProc) prochSOS(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("prochSOS")
	resp := protocol.NewCmdSimple("C"+cmd.GetCmdName()[1:])
	IMEI := cmd.GetInfos()["IMEI"]
	resp.AddArg("1")
	resp.Infos["IMEI"]=IMEI

	if session != nil {
		err := session.Send(libnet.Json(resp))
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}


func (self *ProtoProc) prochupdateSetting(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("prochupdateSetting")
	resp := protocol.NewCmdSimple("C"+cmd.GetCmdName()[1:])
	IMEI := cmd.GetInfos()["IMEI"]
	resp.AddArg("1")
	resp.Infos["IMEI"]=IMEI

	if session != nil {
		err := session.Send(libnet.Json(resp))
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}
