//接入服务器主动连接到推送服务器
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
   心跳检测推送服务器是否存活
*/
func (self *ConnectServer) scanDeadPushClient() {
	timer := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-timer.C:
			log.Info("begin scanDeadClient")
			go func() {
				for ps, pushServerClient := range self.pushServerClientMap {
					self.pushServerClientMutex.Lock()
					if pushServerClient.Alive == false || pushServerClient.Valid == false {
						//心跳没有收到回复，链接作废
						log.Warningf("CloseDeadClient [%s],Alive=%t,Valid=%t.", ps, pushServerClient.Alive, pushServerClient.Valid)
						pushServerClient.Session.Close()
					} else {
						//发送心跳，等待回复
						pushServerClient.Alive = false
						cmd := protocol.NewCmdSimple(protocol.SEND_PING_CMD)
						cmd.AddArg(self.cfg.UUID)
						err := pushServerClient.Session.Send(libnet.Json(cmd))
						if err != nil {
							pushServerClient.Session.Close()
						}
					}
					self.pushServerClientMutex.Unlock()
				}
			}()
		}
	}
}

func (self *ConnectServer) removePushServer(pushServer string) error {
	self.pushServerClientMutex.Lock()
	defer self.pushServerClientMutex.Unlock()

	for index, ps := range self.cfg.PushServerList {
		if ps == pushServer {
			if self.pushServerClientMap[ps] != nil {
				self.pushServerClientMap[ps].Session.Close()
			}
			self.cfg.PushServerList = append(self.cfg.PushServerList[:index], self.cfg.PushServerList[index+1:]...)
		}
	}
	return nil
}

func (self *ConnectServer) addPushServer(pushServer string) error {
	self.pushServerClientMutex.Lock()
	defer self.pushServerClientMutex.Unlock()
	exist := false
	for _, ps := range self.cfg.PushServerList {
		if ps == pushServer {
			exist = true
		}
	}
	if exist == false {
		self.cfg.PushServerList = append(self.cfg.PushServerList, pushServer)
		go self.subscribePushServerChannels()
	}
	return nil
}

/*
   用于反复检测没有连接成功的推送服务器，进行重连
*/
func (self *ConnectServer) subscribePushServerChannels() error {
	log.Info("connect_server start to subscribePushServerChannels")
	self.pushServerClientMutex.Lock()
	defer self.pushServerClientMutex.Unlock()
	for _, ps := range self.cfg.PushServerList {
		if self.pushServerClientMap[ps] != nil {
			//已经创建过链接并且链接正常
			continue
		}
		pushServerClient, err := self.connectPushServer(ps) //发起连接
		if err != nil {
			log.Error(err.Error())
			go self.subscribePushServerChannels()
			continue
		}
		//连接建立成功，开始发送通道订阅
		cmd := protocol.NewCmdSimple(protocol.SUBSCRIBE_CHANNEL_CMD)
		cmd.AddArg(protocol.SYSCTRL_CONNECT_SERVER)
		cmd.AddArg(self.cfg.UUID)

		err = pushServerClient.Send(libnet.Json(cmd))
		if err != nil {
			log.Error(err.Error())
			go self.subscribePushServerChannels()
			continue
		}
		//通道订阅发送成功
		self.pushServerClientMap[ps] = new(pushServerClientState)
		self.pushServerClientMap[ps].Alive = true
		self.pushServerClientMap[ps].Valid = false
		self.pushServerClientMap[ps].ClientSessionNum = 0
		self.pushServerClientMap[ps].Session = pushServerClient

		//开始处理 推送服务器-> 接入服务器 的数据
		go func(ps string) {
			// go self.removePushServer(ps)
			err := self.handlePushServerClient(pushServerClient)
			log.Infof("err=%s", err)
			if err != nil {
				//				if self.pushServerClientMap[ps].Valid == true {
				//					self.pushServerClientNum--
				//					if self.pushServerClientNum == 0 {
				//						self.pushServerClientEmptyMutex.Lock()
				//					}
				//				}
				// self.pushServerClientRWMutex.Lock()
				// defer self.pushServerClientRWMutex.Unlock()
				delete(self.pushServerClientMap, ps)
				log.Info("delete ok")
			}
			go self.subscribePushServerChannels()
		}(ps)
	}
	return nil
}

/*
   处理推送服务器发送过来的数据
*/

func (self *ConnectServer) handlePushServerClient(msc *libnet.Session) error {
	err := msc.Process(func(msg *libnet.InBuffer) error {
		var c protocol.CmdSimple
		ps := msc.Conn().RemoteAddr().String()
		if self.pushServerClientMap[ps] == nil {
			log.Error(ps + " not exist")
			return errors.New(ps + " not exist")
		}
		err := json.Unmarshal(msg.Data, &c)
		if err != nil {
			log.Error("error:", err)
			return err
		}

		log.Infof("c.GetCmdName()=%s\n\n", c.GetCmdName())
		switch c.GetCmdName() {
		case protocol.SUBSCRIBE_CHANNEL_CMD_ACK:
			self.pushServerClientMap[ps].Valid = true
			//			self.pushServerClientNum++
			//			if self.pushServerClientNum == 1 {
			//				self.pushServerClientEmptyMutex.Unlock()
			//			}
		case protocol.PING_CMD_ACK:
			self.pushServerClientMap[ps].Alive = true
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
   连接到推送服务器
*/
func (self *ConnectServer) connectPushServer(ps string) (*libnet.Session, error) {
	client, err := libnet.Dial("tcp", ps)
	if err != nil {
		log.Error(err.Error())
		// panic(err)
	}

	return client, err
}
