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
	"errors"
	"flag"
	"fmt"
	"strconv"

	"github.com/oikomi/FishChatServer/common"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/protocol"
)

var InputConfFile = flag.String("conf_file", "client.json", "input conf file name")

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

func heartBeat(cfg *ClientConfig, msgServerClient *libnet.Session) {
	hb := common.NewHeartBeat("client", msgServerClient, cfg.HeartBeatTime, cfg.Expire, 10)
	hb.Beat()
}

func DisplayCommandList() int {
loopInput:
	fmt.Println("RDA Watch Client Command:")
	fmt.Println("0. Get your topic list")
	fmt.Println("1. Get member list of specific topic")
	fmt.Println("2. Create topic")
	fmt.Println("3. Add someone into the topic, only for topic creator")
	fmt.Println("4. Kick someone out of the topic, only for topic creator")
	fmt.Println("5. Join topic")
	fmt.Println("6. Quit topic")
	fmt.Println("7. Send p2p message")
	fmt.Println("8. Send topic message")
	fmt.Print("Please input number (0~8): ")

	var input string
	var num int
	var err error
	if _, err = fmt.Scanf("%s\n", &input); err != nil {
		log.Error(err.Error())
	}

	if num, err = strconv.Atoi(input); err != nil {
		fmt.Println("Input error")
		goto loopInput
	}
	if num < 0 || num > 8 {
		fmt.Println("Input error")
		goto loopInput
	}
	return num
}

type Client struct {
	clientID   string
	clientType string
	uuid       string
	msAddr     string
	msClient   *libnet.Session
}

func NewClient() *Client {
	return new(Client)
}

// return (ClientID, MsgServerAddr)
func (self *Client) procLoginGateway(cfg *ClientConfig) error {
	fmt.Println("req GatewayServer...")

	gatewayClient, err := libnet.Dial("tcp", cfg.GatewayServer)
	if err != nil {
		panic(err)
	}

	// input client id
	fmt.Println("input my id :")
	var myID string
	var myType string
	var myPass string

	if _, err := fmt.Scanf("%s\n", &myID); err != nil {
		log.Error(err.Error())
	}
	self.clientID = myID

	// input client type
loopType:
	fmt.Println("input my type, D for Device, C for Client :")
	if _, err = fmt.Scanf("%s\n", &myType); err != nil {
		log.Error(err.Error())
	}
	if myType == "D" || myType == "d" {
		myType = protocol.DEV_TYPE_WATCH
	} else if myType == "C" || myType == "c" {
		myType = protocol.DEV_TYPE_CLIENT
	} else {
		goto loopType
	}
	self.clientType = myType

	// input password
	myPass = ""
	if myType == protocol.DEV_TYPE_CLIENT {
		fmt.Println("input my password :")
		if _, err = fmt.Scanf("%s\n", &myPass); err != nil {
			log.Error(err.Error())
		}
	}
	cmd := protocol.NewCmdSimple(protocol.REQ_LOGIN_CMD)
	cmd.AddArg(myID)
	cmd.AddArg(myType)
	cmd.AddArg(myPass)

	err = gatewayClient.Send(libnet.Json(cmd))
	if err != nil {
		log.Error(err.Error())
	}

	var c protocol.CmdSimple
	err = gatewayClient.ProcessOnce(func(msg *libnet.InBuffer) error {
		log.Info(string(msg.Data))
		err = json.Unmarshal(msg.Data, &c)
		if err != nil {
			log.Error("error:", err)
		}
		return nil
	})
	if err != nil {
		log.Error(err.Error())
	}

	gatewayClient.Close()

	fmt.Println("req GatewayServer end...")

	if c.GetArgs()[0] != protocol.RSP_SUCCESS {
		log.Errorf("login gateway error: %s", c.GetArgs()[0])
		return errors.New(c.GetArgs()[0])
	}
	self.uuid = c.GetArgs()[1]
	self.msAddr = c.GetArgs()[2]

	return nil
}

func (self *Client) procLoginServer() error {
	var err error
	var c protocol.Cmd

	self.msClient, err = libnet.Dial("tcp", self.msAddr)
	if err != nil {
		panic(err)
	}

	fmt.Println("req to login msg server...")
	cmd := protocol.NewCmdSimple(protocol.REQ_LOGIN_CMD)
	cmd.AddArg(self.clientID)
	cmd.AddArg(self.uuid)

	err = self.msClient.Send(libnet.Json(cmd))
	if err != nil {
		log.Error(err.Error())
		return err
	}

	err = self.msClient.ProcessOnce(func(msg *libnet.InBuffer) error {
		log.Info(string(msg.Data))
		err = json.Unmarshal(msg.Data, &c)
		if err != nil {
			log.Error("error:", err)
		}
		return nil
	})
	if err != nil {
		log.Error(err.Error())
	}
	if c.GetArgs()[0] != protocol.RSP_SUCCESS {
		log.Errorf("login msgserver error: %s", c.GetArgs()[0])
		return errors.New(c.GetArgs()[0])
	}

	fmt.Println("login msg server SUCCESS")
	return nil
}

func (self *Client) procGetTopicList() error {
	// get topic list
	cmd := protocol.NewCmdSimple(protocol.REQ_GET_TOPIC_LIST_CMD)
	err := self.msClient.Send(libnet.Json(cmd))
	if err != nil {
		log.Error(err.Error())
	}

	return err
}

func (self *Client) procGetTopicListRsp(c *protocol.CmdSimple) error {
	var num int
	var err error

	fmt.Println(c.GetCmdName() + " returns: " + c.GetArgs()[0])
	if c.GetArgs()[0] != protocol.RSP_SUCCESS {
		return errors.New(c.GetArgs()[0])
	}
	if num, err = strconv.Atoi(c.GetArgs()[1]); err != nil {
		fmt.Println(err.Error())
		log.Error(err.Error())
		return err
	}
	fmt.Println("GET_TOPIC_LIST returns (" + c.GetArgs()[1] + "): ")
	index := 0
	for {
		if index == num {
			break
		} else {
			fmt.Println(c.GetArgs()[2+index])
			index++
		}
	}

	return nil
}
func (self *Client) procGetTopicMember() error {
	// get topic list
	cmd := protocol.NewCmdSimple(protocol.REQ_GET_TOPIC_MEMBER_CMD)
	err := self.msClient.Send(libnet.Json(cmd))
	if err != nil {
		log.Error(err.Error())
	}

	return err
}

func (self *Client) procGetTopicMemberRsp(c *protocol.CmdSimple) error {
	var num int
	var err error

	fmt.Println(c.GetCmdName() + " returns: " + c.GetArgs()[0])
	if c.GetArgs()[0] != protocol.RSP_SUCCESS {
		return errors.New(c.GetArgs()[0])
	}
	if num, err = strconv.Atoi(c.GetArgs()[1]); err != nil {
		fmt.Println(err.Error())
		log.Error(err.Error())
		return err
	}
	fmt.Println("GET_TOPIC_MEMBER returns (" + c.GetArgs()[1] + "): ")
	index := 0
	for {
		if index == num {
			break
		} else {
			fmt.Println("ID=" + c.GetArgs()[2+2*index] + "\t\t\t, Name=" + c.GetArgs()[2+2*index+1])
			index++
		}
	}

	return nil
}

func (self *Client) procCreateTopic() error {
	// CREATE TOPIC
	var input string
	fmt.Println("want to create a topic (y/n) :")
	if _, err := fmt.Scanf("%s\n", &input); err != nil {
		log.Error(err.Error())
	}
	if input == "y" {
		cmd := protocol.NewCmdSimple(protocol.REQ_CREATE_TOPIC_CMD)
		fmt.Println("CREATE_TOPIC_CMD | input topic name :")
		if _, err := fmt.Scanf("%s\n", &input); err != nil {
			log.Error(err.Error())
		}
		cmd.AddArg(input)

		fmt.Println("CREATE_TOPIC_CMD | input alias name :")
		if _, err := fmt.Scanf("%s\n", &input); err != nil {
			log.Error(err.Error())
		}
		cmd.AddArg(input)

		err := self.msClient.Send(libnet.Json(cmd))
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}

func (self *Client) procCreateTopicRsp(c *protocol.CmdSimple) error {
	fmt.Println(c.GetCmdName() + " returns: " + c.GetArgs()[0])
	if c.GetArgs()[0] != protocol.RSP_SUCCESS {
		return errors.New(c.GetArgs()[0])
	}
	return nil
}

func (self *Client) procJoinTopic() error {
	var input string
	fmt.Println("want to join a topic (y/n) :")
	if _, err := fmt.Scanf("%s\n", &input); err != nil {
		log.Error(err.Error())
		return err
	}
	if input == "y" {
		cmd := protocol.NewCmdSimple(protocol.REQ_JOIN_TOPIC_CMD)

		fmt.Println("input topic name :")
		if _, err := fmt.Scanf("%s\n", &input); err != nil {
			fmt.Errorf(err.Error())
			return err
		}
		cmd.AddArg(input)

		fmt.Println("input alias name :")
		if _, err := fmt.Scanf("%s\n", &input); err != nil {
			fmt.Errorf(err.Error())
			return err
		}
		cmd.AddArg(input)

		err := self.msClient.Send(libnet.Json(cmd))
		if err != nil {
			fmt.Errorf(err.Error())
			return err
			log.Error(err.Error())
		}
	}
	return nil
}

func (self *Client) procJoinTopicRsp(c *protocol.CmdSimple) error {
	fmt.Println(c.GetCmdName() + " returns: " + c.GetArgs()[0])
	if c.GetArgs()[0] != protocol.RSP_SUCCESS {
		return errors.New(c.GetArgs()[0])
	}
	return nil
}

func (self *Client) procQuitTopic() error {
	var input string
	fmt.Println("want to quit a topic (y/n) :")
	if _, err := fmt.Scanf("%s\n", &input); err != nil {
		log.Error(err.Error())
		return err
	}
	if input == "y" {
		cmd := protocol.NewCmdSimple(protocol.REQ_QUIT_TOPIC_CMD)

		fmt.Println("input topic name :")
		if _, err := fmt.Scanf("%s\n", &input); err != nil {
			fmt.Errorf(err.Error())
			return err
		}
		cmd.AddArg(input)

		err := self.msClient.Send(libnet.Json(cmd))
		if err != nil {
			fmt.Errorf(err.Error())
			return err
			log.Error(err.Error())
		}
	}
	return nil
}

func (self *Client) procQuitTopicRsp(c *protocol.CmdSimple) error {
	fmt.Println(c.GetCmdName() + " returns: " + c.GetArgs()[0])
	if c.GetArgs()[0] != protocol.RSP_SUCCESS {
		return errors.New(c.GetArgs()[0])
	}
	return nil
}

func (self *Client) procAdd2Topic() error {
	var input string
	fmt.Println("want to add members into a topic (y/n) :")
	if _, err := fmt.Scanf("%s\n", &input); err != nil {
		log.Error(err.Error())
		return err
	}
	if input == "y" {
		cmd := protocol.NewCmdSimple(protocol.REQ_ADD_2_TOPIC_CMD)

		fmt.Println("input topic name :")
		if _, err := fmt.Scanf("%s\n", &input); err != nil {
			fmt.Errorf(err.Error())
			return err
		}
		cmd.AddArg(input)

		fmt.Println("input member ID :")
		if _, err := fmt.Scanf("%s\n", &input); err != nil {
			fmt.Errorf(err.Error())
			return err
		}
		cmd.AddArg(input)

		fmt.Println("input member name :")
		if _, err := fmt.Scanf("%s\n", &input); err != nil {
			fmt.Errorf(err.Error())
			return err
		}
		cmd.AddArg(input)

		err := self.msClient.Send(libnet.Json(cmd))
		if err != nil {
			fmt.Errorf(err.Error())
			return err
			log.Error(err.Error())
		}
	}
	return nil
}

func (self *Client) procAdd2TopicRsp(c *protocol.CmdSimple) error {
	fmt.Println(c.GetCmdName() + " returns: " + c.GetArgs()[0])
	if c.GetArgs()[0] != protocol.RSP_SUCCESS {
		return errors.New(c.GetArgs()[0])
	}
	return nil
}

func (self *Client) procKickTopic() error {
	var input string
	fmt.Println("want to kick member out of a topic (y/n) :")
	if _, err := fmt.Scanf("%s\n", &input); err != nil {
		log.Error(err.Error())
		return err
	}
	if input == "y" {
		cmd := protocol.NewCmdSimple(protocol.REQ_KICK_TOPIC_CMD)

		fmt.Println("input topic name :")
		if _, err := fmt.Scanf("%s\n", &input); err != nil {
			fmt.Errorf(err.Error())
			return err
		}
		cmd.AddArg(input)

		fmt.Println("input member ID :")
		if _, err := fmt.Scanf("%s\n", &input); err != nil {
			fmt.Errorf(err.Error())
			return err
		}
		cmd.AddArg(input)

		err := self.msClient.Send(libnet.Json(cmd))
		if err != nil {
			fmt.Errorf(err.Error())
			return err
			log.Error(err.Error())
		}
	}
	return nil
}

func (self *Client) procKickTopicRsp(c *protocol.CmdSimple) error {
	fmt.Println(c.GetCmdName() + " returns: " + c.GetArgs()[0])
	if c.GetArgs()[0] != protocol.RSP_SUCCESS {
		return errors.New(c.GetArgs()[0])
	}
	return nil
}

func (self *Client) procSendP2PMsg() error {
	var input string

	cmd := protocol.NewCmdSimple(protocol.REQ_SEND_P2P_MSG_CMD)

	fmt.Println("send the id you want to talk :")
	if _, err := fmt.Scanf("%s\n", &input); err != nil {
		log.Error(err.Error())
	}

	cmd.AddArg(input)

	fmt.Println("input msg :")
	if _, err := fmt.Scanf("%s\n", &input); err != nil {
		log.Error(err.Error())
	}

	cmd.AddArg(input)

	err := self.msClient.Send(libnet.Json(cmd))
	if err != nil {
		log.Error(err.Error())
	}
	return nil
}

func (self *Client) procSendP2PMsgRsp(c *protocol.CmdSimple) error {
	fmt.Println(c.GetCmdName() + " returns: " + c.GetArgs()[0])
	if c.GetArgs()[0] == protocol.RSP_SUCCESS {
		fmt.Println("uuid: " + c.GetArgs()[1])
	}
	return nil
}

func (self *Client) procSendP2PMsgReq(c *protocol.CmdSimple) error {
	fmt.Println(c.GetArgs()[1] + "  says : " + c.GetArgs()[0])
	if len(c.GetArgs()) >= 3 {
		cmd := protocol.NewCmdSimple(protocol.IND_ACK_P2P_MSG_CMD)
		cmd.AddArg(c.GetArgs()[2])
		cmd.AddArg(protocol.P2P_ACK_READ)

		err := self.msClient.Send(libnet.Json(cmd))
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}

func (self *Client) procSendTopicMsg() error {
	var input string

	cmd := protocol.NewCmdSimple(protocol.REQ_SEND_TOPIC_MSG_CMD)

	fmt.Println("input msg :")
	if _, err := fmt.Scanf("%s\n", &input); err != nil {
		log.Error(err.Error())
	}
	cmd.AddArg(input)

	fmt.Println("input the topic you want to talk :")
	if _, err := fmt.Scanf("%s\n", &input); err != nil {
		log.Error(err.Error())
	}
	cmd.AddArg(input)

	err := self.msClient.Send(libnet.Json(cmd))
	if err != nil {
		log.Error(err.Error())
	}
	return nil
}

func (self *Client) procSendTopicMsgRsp(c *protocol.CmdSimple) error {
	fmt.Println(c.GetCmdName() + " returns: " + c.GetArgs()[0])
	return nil
}

func (self *Client) procSendTopicMsgReq(c *protocol.CmdSimple) error {
	msg := c.GetArgs()[0]
	topicName := c.GetArgs()[1]
	fromID := c.GetArgs()[2]
	fromType := c.GetArgs()[3]
	fmt.Println("Topic message received :")
	fmt.Println("    TopicName :" + topicName)
	fmt.Println("    FromID    :" + fromID)
	fmt.Println("    FromType  :" + fromType)
	fmt.Println("    Message   :" + msg)
	return nil
}

func main() {
	flag.Parse()
	cfg, err := LoadConfig(*InputConfFile)
	if err != nil {
		log.Error(err.Error())
		return
	}

	client := NewClient()
	err = client.procLoginGateway(cfg)
	if err != nil {
		panic(err)
	}

	err = client.procLoginServer()
	if err != nil {
		panic(err)
	}

	go heartBeat(cfg, client.msClient)

	var c protocol.CmdSimple
	go client.msClient.Process(func(msg *libnet.InBuffer) error {
		log.Info(string(msg.Data))
		err = json.Unmarshal(msg.Data, &c)
		if err != nil {
			log.Error("error:", err)
		}

		fmt.Println("msg received is : " + c.GetCmdName())
		switch c.GetCmdName() {

		case protocol.RSP_GET_TOPIC_LIST_CMD:
			client.procGetTopicListRsp(&c)

		case protocol.RSP_GET_TOPIC_MEMBER_CMD:
			client.procGetTopicMemberRsp(&c)

		case protocol.RSP_CREATE_TOPIC_CMD:
			client.procCreateTopicRsp(&c)

		case protocol.RSP_JOIN_TOPIC_CMD:
			client.procJoinTopicRsp(&c)

		case protocol.RSP_QUIT_TOPIC_CMD:
			client.procQuitTopicRsp(&c)

		case protocol.RSP_ADD_2_TOPIC_CMD:
			client.procAdd2TopicRsp(&c)

		case protocol.RSP_KICK_TOPIC_CMD:
			client.procKickTopicRsp(&c)

		case protocol.RSP_SEND_P2P_MSG_CMD:
			client.procSendP2PMsgRsp(&c)

		case protocol.IND_ACK_P2P_MSG_CMD:
			fmt.Println("msg sent [uuid=" + c.GetArgs()[0] + "] status: " + c.GetArgs()[1])

		case protocol.REQ_SEND_P2P_MSG_CMD:
			client.procSendP2PMsgReq(&c)

		case protocol.REQ_SEND_TOPIC_MSG_CMD:
			client.procSendTopicMsgReq(&c)
		}

		return nil
	})

	for {
		num := DisplayCommandList()
		/*
			fmt.Println("0. Get your topic list")
			fmt.Println("1. Get member list of specific topic")
			fmt.Println("2. Create topic")
			fmt.Println("3. Add someone into the topic, only for topic creator")
			fmt.Println("4. Kick someone out of the topic, only for topic creator")
			fmt.Println("5. Join topic")
			fmt.Println("6. Quit topic")
			fmt.Println("7. Send p2p message")
			fmt.Println("8. Send topic message")
		*/
		switch num {
		case 0:
			client.procGetTopicList()
		case 1:
			client.procGetTopicMember()
		case 2:
			client.procCreateTopic()
		case 3:
			client.procAdd2Topic()
		case 4:
			client.procKickTopic()
		case 5:
			client.procJoinTopic()
		case 6:
			client.procQuitTopic()
		case 7:
			client.procSendP2PMsg()
		case 8:
			client.procSendTopicMsg()
		}
	}

	defer client.msClient.Close()

	// msgServerClient.Process(func(msg *libnet.InBuffer) error {
	// 	log.Info(string(msg.Data))
	// 	return nil
	// })

	log.Flush()
}
