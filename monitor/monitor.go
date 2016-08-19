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
	//	"flag"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/oikomi/FishChatServer/monitor/service"
	//	"github.com/oikomi/FishChatServer/log"
	//	_ "github.com/oikomi/FishChatServer/monitor/docs"
	//	"github.com/oikomi/FishChatServer/monitor/controllers"
	//	_ "github.com/oikomi/FishChatServer/monitor/routers"
	_ "github.com/oikomi/FishChatServer/monitor/routers"
)

const VERSION string = "0.10"

func version() {
	fmt.Printf("monitor version %s Copyright (c) 2014-2015 Harold Miao (miaohong@miaohong.org)  \n", VERSION)
}

func main() {
	version()
	service.Init()
	beego.Run()
}
