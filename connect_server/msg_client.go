//接入服务器主动连接到消息服务器
package main

import (
    "encoding/json"
    "fmt"
    "time"

    "github.com/oikomi/FishChatServer/libnet"
    "github.com/oikomi/FishChatServer/log"
    "github.com/oikomi/FishChatServer/protocol"
)
/*
    心跳检测消息服务器是否存活
*/
func (self *ConnectServer) scanDeadClient() {
    log.Info("scanDeadClient")
    timer := time.NewTicker(5 * time.Second)
    for {
        select {
        case <-timer.C:
            log.Info("scanDeadClient timeout")
            go func() {
                for _, ms := range self.cfg.MsgServerList {
                    self.msgServerClientMutex.Lock()
                    fmt.Printf(ms)
                    if self.msgServerClientMap[ms].Alive == false {
                        //心跳没有收到回复，链接作废
                        fmt.Printf("Close:%s",ms)
                        self.msgServerClientMap[ms].Session.Close()
                        // delete(self.msgServerClientMap, ms)
                    } else {
                        //发送心跳，等待回复
                        self.msgServerClientMap[ms].Alive = false
                        cmd := protocol.NewCmdSimple(protocol.SUBSCRIBE_CHANNEL_CMD)
                        cmd.AddArg(protocol.SYSCTRL_CONNECT_SERVER)
                        cmd.AddArg(self.cfg.UUID)
                        err := self.msgServerClientMap[ms].Session.Send(libnet.Json(cmd))
                        if err != nil {
                            self.msgServerClientMap[ms].Session.Close()
                        }
                    }
                    self.msgServerClientMutex.Unlock()
                }
            }()
        }
    }
}


/*
    用于反复检测没有连接成功的消息服务器，进行重连
*/
func (self *ConnectServer)subscribeChannels() error {
    log.Info("connect_server start to subscribeChannels")
    self.msgServerClientMutex.Lock()
    defer self.msgServerClientMutex.Unlock()
    for _, ms := range self.cfg.MsgServerList {
        if self.msgServerClientMap[ms] != nil  {
            log.Info("self.msgServerClientMap[ms] != nil")
            continue
        }
        msgServerClient, err := self.connectMsgServer(ms)
        if err != nil {
            log.Error(err.Error())
            go self.subscribeChannels()
            continue
        }
        cmd := protocol.NewCmdSimple(protocol.SUBSCRIBE_CHANNEL_CMD)
        cmd.AddArg(protocol.SYSCTRL_CONNECT_SERVER)
        cmd.AddArg(self.cfg.UUID)
        
        err = msgServerClient.Send(libnet.Json(cmd))
        if err != nil {
            log.Error(err.Error())
            go self.subscribeChannels()
            continue
        }
        self.msgServerClientMap[ms]=new(msgServerClientState)
        self.msgServerClientMap[ms].Alive = true
        self.msgServerClientMap[ms].Session = msgServerClient
        go func(ms string) {
            err := self.handleMsgServerClient(msgServerClient)
            log.Info(" go msgServerClientMap[%s],err=%s",ms,err)
            if err !=nil {
                delete(self.msgServerClientMap, ms)
            }
            self.subscribeChannels()
        }(ms)
    }
    return nil
}


/*
    处理消息服务器发送过来的数据
*/

func (self *ConnectServer)handleMsgServerClient(msc *libnet.Session) error {
    err := msc.Process(func(msg *libnet.InBuffer) error {
        //log.Info("msg_server", msc.Conn().RemoteAddr().String()," say: ", string(msg.Data))
        var c protocol.CmdMonitor
        
        err := json.Unmarshal(msg.Data, &c)
        if err != nil {
            log.Error("error:", err)
            return err
        }

        return nil
    })
    fmt.Printf("handleMsgServerClient,err=%s\n\n",err)
    return err
}

/*
    连接到消息服务器
*/
func (self *ConnectServer)connectMsgServer(ms string) (*libnet.Session, error) {
    client, err := libnet.Dial("tcp", ms)
    if err != nil {
        log.Error(err.Error())
        // panic(err)
    }

    return client, err
}
