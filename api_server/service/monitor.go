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

package service

import (
	"flag"
	"fmt"

	"github.com/oikomi/FishChatServer/log"
	//	_ "github.com/oikomi/FishChatServer/monitor/docs"
	//	"github.com/oikomi/FishChatServer/monitor/controllers"
	//	_ "github.com/oikomi/FishChatServer/monitor/routers"
)

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
const char* build_time(void) {
	static const char* psz_build_time = "["__DATE__ " " __TIME__ "]";
	return psz_build_time;
}
*/
//import "C"

//var (
//	buildTime = C.GoString(C.build_time())
//)

//func BuildTime() string {
//	return buildTime
//}

const VERSION string = "0.10"

func init() {
	flag.Set("alsologtostderr", "false")
	flag.Set("log_dir", "false")
}

func version() {
	fmt.Printf("monitor version %s Copyright (c) 2014-2015 Harold Miao (miaohong@miaohong.org)  \n", VERSION)
}

var InputConfFile = flag.String("conf_file", "monitor.json", "input conf file name")
var m *Monitor

func GetServer() *Monitor {
	return m
}

func Init() {
	version()
	//	fmt.Printf("built on %s\n", BuildTime())
	flag.Parse()
	cfg := NewMonitorConfig(*InputConfFile)
	err := cfg.LoadConfig()
	if err != nil {
		log.Error(err.Error())
		return
	}

	m = NewMonitor(cfg)

	go m.scanDeadClient() //清理无用消息服务器

	m.subscribeChannels()

}
