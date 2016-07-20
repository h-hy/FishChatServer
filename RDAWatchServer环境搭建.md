# RDAWatchServer环境搭建 #

## git环境搭建 ##
参见文档《HOWTO--git&github.md》

## Golang环境搭建 ##
参见文档《golang环境搭建》
将GOPATH设置为$HOME/RDAWatchServer

## redis安装 ##
FishChatServer通过Redis(http://redis.io)做cache，最新稳定版本为3.2.0  
参见文档《redis》

## mongoDB安装 ##
FishChatServer使用MongoDB(http://www.mongodb.org/)做持久化存储，最新稳定版本为3.2.6  
参见文档《mongodb》

## RDAWatchServer代码获取 ##
<pre>
mkdir -p RDAWatchServer/src
cd RDAWatchServer/src
mkdir -p gopkg.in/mgo.v2
mkdir -p github.com/astaxie/beego
mkdir -p github.com/garyburd/redigo
mkdir -p github.com/oikomi/FishChatServer

//然后通过go get获取代码
go get gopkg.in/mgo.v2 # MongoDB驱动
go get github.com/astaxie/beego #web监控使用的beego框架
go get github.com/garyburd/redigo #redis驱动
go get github.com/alvin921/FishChatServer #server代码

//但是失败，原因未知，只好使用git clone获取代码
git clone https://gopkg.in/mgo.v2 gopkg.in/mgo.v2
git clone https://github.com/astaxie/beego github.com/astaxie/beego
git clone https://github.com//garyburd/redigo github.com//garyburd/redigo
git clone https://github.com/alvin921/FishChatServer github.com/oikomi/FishChatServer
</pre>

## 编译及安装 ##
写了如下三个脚本（windows7下，请将cmd.exe和git-bash.exe设置为管理员身份运行）：  
wincmd.bat： 切换到windows cmd命令行窗口  
<pre>
@echo off
@title RDAWatchServer
c:\Windows\system32\cmd.exe 
</pre>

gitcmd.bat： 切换到git命令行窗口
<pre>
@echo off
@title RDAWatchServer
@for /f %%i in ('cd') do set PWD=%%i
c:\Windows\system32\cmd.exe /c ""C:\Program Files\Git\git-bash.exe" --cd=%PWD%"
</pre>

r.sh：在gitcmd命令行运行，主要用于编译和启动/停止服务

<pre>
Usage:
./r.sh clean|build|start|stop ...
./r.sh clean
    clean exe files of gateway/msg_server/manager/router/monitor/client
./r.sh build nil|server|gateway|msg_server|manager|router|monitor|client
    nil|server: means to build all: gateway/msg_erver/manager/router/monitor/client
./r.sh start nil|server|redis|mongo
    nil|server: means to start all: msg_erver/gateway/manager/router/monitor
./r.sh stop  nil|server|redis|mongo
    nil|server: means to stop all: msg_erver/gateway/manager/router/monitor
</pre>

## 服务器部署 ##
FishChatServer采用分布式可伸缩部署方式(各类服务器角色都可以动态增减)：

*   gateway一台
*   msg_server两台
*   router一台
*   manager一台
*   monitor一台

如果没有多机条件，可以单机部署。


**NOTE:**  必须先修改各文件夹下面的json配置文件配置服务器参数

**NOTE:** gateway、router、manager和monitor一定要在msg_server之后启动，因为他们都订阅了msg_server的channel 

**NOTE:** 

## 测试 ##

<pre>
./r.sh build client
client/client
</pre>
