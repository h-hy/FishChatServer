//接入服务器主动连接到消息服务器
package main

import (
	"encoding/json"
	"errors"
	_ "fmt"
	"time"

	"github.com/oikomi/FishChatServer/connect_libnet"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/protocol"
)

/*
   心跳检测消息服务器是否存活
*/
func (self *ConnectServer) scanDeadMsgClient() {
	timer := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-timer.C:
			log.Info("begin scanDeadClient")
			go func() {
				for ms, msgServerClient := range self.msgServerClientMap {
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

func (self *ConnectServer) removeMsgServer(msgServer string) error {
	self.msgServerClientMutex.Lock()
	defer self.msgServerClientMutex.Unlock()

	for index, ms := range self.cfg.MsgServerList {
		if ms == msgServer {
			if self.msgServerClientMap[ms] != nil {
				self.msgServerClientMap[ms].Session.Close()
			}
			self.cfg.MsgServerList = append(self.cfg.MsgServerList[:index], self.cfg.MsgServerList[index+1:]...)
		}
	}
	return nil
}

func (self *ConnectServer) addMsgServer(msgServer string) error {
	self.msgServerClientMutex.Lock()
	defer self.msgServerClientMutex.Unlock()
	exist := false
	for _, ms := range self.cfg.MsgServerList {
		if ms == msgServer {
			exist = true
		}
	}
	if exist == false {
		self.cfg.MsgServerList = append(self.cfg.MsgServerList, msgServer)
		go self.subscribeChannels()
	}
	return nil
}

/*
   用于反复检测没有连接成功的消息服务器，进行重连
*/
func (self *ConnectServer) subscribeChannels() error {
	log.Info("connect_server start to subscribeChannels")
	self.msgServerClientMutex.Lock()
	defer self.msgServerClientMutex.Unlock()
	for _, ms := range self.cfg.MsgServerList {
		if self.msgServerClientMap[ms] != nil {
			//已经创建过链接并且链接正常
			continue
		}
		msgServerClient, err := self.connectMsgServer(ms) //发起连接
		if err != nil {
			log.Error(err.Error())
			go self.subscribeChannels()
			continue
		}
		//连接建立成功，开始发送通道订阅
		cmd := protocol.NewCmdSimple(protocol.SUBSCRIBE_CHANNEL_CMD)
		cmd.AddArg(protocol.SYSCTRL_CONNECT_SERVER)
		cmd.AddArg(self.cfg.UUID)

		err = msgServerClient.Send(libnet.Json(cmd))
		if err != nil {
			log.Error(err.Error())
			go self.subscribeChannels()
			continue
		}
		//通道订阅发送成功
		self.msgServerClientMap[ms] = new(msgServerClientState)
		self.msgServerClientMap[ms].Alive = true
		self.msgServerClientMap[ms].Valid = false
		self.msgServerClientMap[ms].ClientSessionNum = 0
		self.msgServerClientMap[ms].Session = msgServerClient

		//开始处理 消息服务器-> 接入服务器 的数据
		go func(ms string) {
			// go self.removeMsgServer(ms)
			err := self.handleMsgServerClient(msgServerClient)
			log.Infof("err=%s", err)
			if err != nil {
				if self.msgServerClientMap[ms].Valid == true {
					self.msgServerClientNum--
					if self.msgServerClientNum == 0 {
						self.msgServerClientEmptyMutex.Lock()
					}
				}
				// self.msgServerClientRWMutex.Lock()
				// defer self.msgServerClientRWMutex.Unlock()
				delete(self.msgServerClientMap, ms)
				log.Info("delete ok")
			}
			go self.subscribeChannels()
		}(ms)
	}
	return nil
}

/*
   处理消息服务器发送过来的数据
*/

func (self *ConnectServer) handleMsgServerClient(msc *libnet.Session) error {
	err := msc.Process(func(msg *libnet.InBuffer) error {
		var c protocol.CmdSimple
		ms := msc.Conn().RemoteAddr().String()
		if self.msgServerClientMap[ms] == nil {
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
			self.msgServerClientMap[ms].Valid = true
			self.msgServerClientNum++
			if self.msgServerClientNum == 1 {
				self.msgServerClientEmptyMutex.Unlock()
			}
		case protocol.PING_CMD_ACK:
			self.msgServerClientMap[ms].Alive = true
		default:
			if c.GetInfos()["IMEI"] != "" {
				IMEI := c.GetInfos()["IMEI"]
				if self.sessions[IMEI] != nil {
					self.sessions[IMEI].Send(connect_libnet.Bytes([]byte(c.GetDatas())))
				}
			}
		}
		return nil
	})
	return err
}

/*
   连接到消息服务器
*/
func (self *ConnectServer) connectMsgServer(ms string) (*libnet.Session, error) {
	client, err := libnet.Dial("tcp", ms)
	if err != nil {
		log.Error(err.Error())
		// panic(err)
	}

	return client, err
}
