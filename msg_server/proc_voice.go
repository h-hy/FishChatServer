package main

import (
	"encoding/base64"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
	// "bytes"

	// "github.com/oikomi/FishChatServer/common"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/models"
	"github.com/oikomi/FishChatServer/protocol"
)

func readFile(pathFilename string, sizePreGroup, groupId int) (int, []byte, error) {

	data := make([]byte, sizePreGroup)
	file, err := os.Open(pathFilename) // For read access.
	defer file.Close()
	if err != nil {
		log.Info("Open File Error", err)
		return 0, data, err
	}
	readed, err := file.ReadAt(data, int64(groupId*sizePreGroup))
	if err == io.EOF && readed > 0 {
		return readed, data, nil
	} else if err != nil {
		log.Info("Read File Error", err)
		return 0, data, err
	}
	return readed, data, nil
}

func writeFile(pathFilename string, data []byte, off int64) (int, error) {
	//	file, err := os.Create(pathFilename) // For read access.

	file, err := os.OpenFile(pathFilename, os.O_RDWR, 0666) // For read access.
	defer file.Close()
	if os.IsNotExist(err) {
		file, err = os.Create(pathFilename) //创建文件
	} else if err != nil {
		log.Info("Create File Error", err)
		return 0, err
	}
	writed, err := file.WriteAt(data, off)
	if err != nil {
		log.Info("Write File Error", err)
		return 0, err
	}
	return writed, nil
}

func (self *ProtoProc) procVoiceDown(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("procVoiceDown")
	IMEI := cmd.GetInfos()["IMEI"]
	newCmd := protocol.NewCmdSimple("D" + cmd.GetCmdName()[1:])
	newCmd.Infos["IMEI"] = IMEI

	log.Info("len(cmd.GetArgs())=", len(cmd.GetArgs()))
	if len(cmd.GetArgs()) < 3 {
		return nil
	}
	id := cmd.GetArgs()[0]
	step := cmd.GetArgs()[1]
	if step == "0" {
		//声明结果
		if len(cmd.GetArgs()) < 2 {
			log.Info("Transfer Step 0 Fail A")
			return nil
		}
		result := cmd.GetArgs()[2]
		if result != "1" {
			log.Info("Transfer Step 0 Fail")
			return nil
		}
		//传送第一组开始
		voiceCacheData, err := self.msgServer.voiceCache.Get("DOWN", id)
		if err != nil {
			log.Info("Transfer Step 0 Error", err)
			return nil
		}
		//文件信息读取完毕
		pathFilename := voiceCacheData.PathFilename
		readed, data, err := readFile(pathFilename, 1024, 0)
		if err != nil {
			log.Info("Transfer Step 0 read File Error", err)
			return nil
		}
		newCmd.AddArg(id)
		newCmd.AddArg("1")
		newCmd.AddArg("1")
		newCmd.AddArg(strconv.Itoa(readed))
		newCmd.AddArg(string(data[:readed]))
		if session != nil {
			if err := session.Send(libnet.Json(newCmd)); err != nil {
				log.Error(err.Error())
			}
		}
		return nil
	} else if step == "1" {
		if len(cmd.GetArgs()) < 4 {
			log.Info("Transfer Step 1 Fail A")
			return nil
		}
		//		groupId := cmd.GetArgs()[2]

		if result := cmd.GetArgs()[3]; result != "1" {
			log.Info("Transfer Step 1 Fail")
			return nil
		}
		if len(cmd.GetArgs()) < 5 {
			log.Info("Transfer Step 1 Fail B")
			return nil
		}
		nextGroupId, err := strconv.Atoi(cmd.GetArgs()[4])
		if err != nil {
			log.Info("Transfer Step 1 ,nextGroupId Error", err)
			return nil
		}
		//传送nextGroupId组
		voiceCacheData, err := self.msgServer.voiceCache.Get("DOWN", id)
		if err != nil {
			log.Info("Transfer Step 1 Error", err)
			return nil
		}
		//文件信息读取完毕
		if 1024*nextGroupId >= voiceCacheData.Size {
			if 1024*(nextGroupId-1) >= voiceCacheData.Size {
				log.Info("nextGroupId Error")
				return nil
			}
			//已经传送完毕！
			newCmd.AddArg(id)
			newCmd.AddArg("2")
			if session != nil {
				if err := session.Send(libnet.Json(newCmd)); err != nil {
					log.Error(err.Error())
				}
			}
			return nil
		}
		pathFilename := voiceCacheData.PathFilename
		readed, data, err := readFile(pathFilename, 1024, nextGroupId)
		if err != nil {
			log.Info("Transfer Step 1 read File Error", err)
			return nil
		}
		newCmd.AddArg(id)
		newCmd.AddArg("1")
		newCmd.AddArg(cmd.GetArgs()[4]) //nextGroupId
		newCmd.AddArg(strconv.Itoa(readed))
		newCmd.AddArg(base64.StdEncoding.EncodeToString(data[:readed]))

		if session != nil {
			if err := session.Send(libnet.Json(newCmd)); err != nil {
				log.Error(err.Error())
			}
		}
		return nil

	} else if step == "2" {
		result := cmd.GetArgs()[2]
		if result != "1" {
			log.Info("Transfer Step 2 Fail")
			return nil
		}
	}
	return nil
}

/*
 *    procVoiceUp  处理文件上行
 *    Huang Haoyan   2016.08
 */
func (self *ProtoProc) procVoiceUp(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("procVoiceUp")
	IMEI := cmd.GetInfos()["IMEI"]
	newCmd := protocol.NewCmdSimple("C" + cmd.GetCmdName()[1:])
	newCmd.Infos["IMEI"] = IMEI

	log.Info("len(cmd.GetArgs())=", len(cmd.GetArgs()))
	if len(cmd.GetArgs()) < 3 {
		return nil
	}
	kid := cmd.GetArgs()[0]
	id, err := strconv.Atoi(cmd.GetArgs()[1])
	if err != nil {
		log.Infof("procVoiceUp, Step 0 Fail, id [%s] error.", cmd.GetArgs()[1], err)
		return nil
	}
	step := cmd.GetArgs()[2]
	newCmd.AddArg(kid)              //kid
	newCmd.AddArg(cmd.GetArgs()[1]) //id
	newCmd.AddArg(step)

	if step == "0" {
		//声明结果
		fail := false
		if len(cmd.GetArgs()) < 5 {
			log.Info("procVoiceUp ,Step 0 Fail for length of Args.")
			fail = true
		}
		format := cmd.GetArgs()[3]
		size, err := strconv.Atoi(cmd.GetArgs()[4])
		if !fail && err != nil {
			log.Infof("procVoiceUp, Step 0 Fail, size [%s] error.", cmd.GetArgs()[4], err)
			fail = true
		}
		if !fail && format != "amr" {
			log.Infof("Transfer Step 0 Fail,format [%s] not support.", format)
			fail = true
		}
		if !fail {
			//创建缓存
			timestr := time.Now().Format("2006_01_02_15_04_05")
			filename := IMEI + "_" + timestr + "_" + string(Krand(5, KC_RAND_KIND_LOWER)) + ".amr"
			UpAMRSaveDir, _ := filepath.Abs(self.msgServer.cfg.VoiceUpSaveDir)
			log.Info(UpAMRSaveDir + filename)
			voiceCacheData := self.msgServer.voiceCache.NewVoiceCacheData("UP", id, "", filename, UpAMRSaveDir+filename, format, size)
			self.msgServer.voiceCache.Set(voiceCacheData)
			//创建缓存完毕   REPLY : KID#ID#0#1#
		}
		if !fail {
			newCmd.AddArg("1")
		} else {
			newCmd.AddArg("2")
		}

		if session != nil {
			if err := session.Send(libnet.Json(newCmd)); err != nil {
				log.Error(err.Error())
			}
		}
	} else if step == "1" {
		//上传数据 U36#KID#ID#1#groupid#size#group data
		fail := false
		if len(cmd.GetArgs()) != 6 {
			log.Info("Transfer Step 1 Fail A")
			fail = true
		}
		voiceCacheData, err := self.msgServer.voiceCache.Get("UP", cmd.GetArgs()[1]) //id
		if err != nil {
			log.Info("Transfer Step 1 ,id Error", err)
			fail = true
		}

		groupId, err := strconv.Atoi(cmd.GetArgs()[3])
		if err != nil {
			log.Info("Transfer Step 1 Fail,groupId error.", groupId, err)
			fail = true
		}
		size, err := strconv.Atoi(cmd.GetArgs()[4])
		if err != nil {
			log.Info("Transfer Step 1 Fail,size error.", size, err)
			fail = true
		}

		if !fail && voiceCacheData.NowGroup < groupId {
			log.Info("Transfer Step 1 ,groupId Error", groupId, voiceCacheData.NowGroup)
			fail = true
		}
		if !fail && voiceCacheData.NowGroup == groupId {
			if voiceCacheData.NowSize+size > voiceCacheData.Size {
				log.Info("Transfer Step 1 ERROR ,NowSize+size max than voiceCacheData.Size", voiceCacheData.Size, size)
				fail = true
			}
			//ok write it
			data, err := base64.StdEncoding.DecodeString(cmd.GetArgs()[5])
			if err != nil {
				log.Info("Transfer Step 1 Fail,DecodeString error.")
				fail = true
			}
			if size != len(data) {
				log.Info("Transfer Step 0 Fail,size error B.", size, len(data))
				fail = true
			}
			if !fail {
				writeFile(voiceCacheData.PathFilename, data, int64(voiceCacheData.NowSize))
				voiceCacheData.NowGroup++
				voiceCacheData.NowSize += size
				self.msgServer.voiceCache.Set(voiceCacheData)
			}
		}
		newCmd.AddArg(cmd.GetArgs()[3]) //groupId
		if !fail {
			//C36#KID #ID#1#group_id#1#next group id#
			newCmd.AddArg("1")
			newCmd.AddArg(strconv.Itoa(voiceCacheData.NowGroup))
		} else {
			newCmd.AddArg("3")
		}

		if session != nil {
			if err := session.Send(libnet.Json(newCmd)); err != nil {
				log.Error(err.Error())
			}
		}
		return nil

	} else if step == "2" {
		//上传数据完毕 U36#KID#ID#1#groupid#size#group data
		fail := false
		voiceCacheData, err := self.msgServer.voiceCache.Get("UP", cmd.GetArgs()[1]) //id
		if err != nil {
			log.Info("Transfer Step 1 ,id Error", err)
			fail = true
		}
		if !fail && voiceCacheData.NowSize != voiceCacheData.Size {
			log.Info("Transfer Step 2 ,NowSize Error", voiceCacheData.Size, voiceCacheData.NowSize)
			fail = true
			return nil
		}
		if !fail {
			amrURIPrefix := self.msgServer.cfg.AmrURIPrefix
			_, err := models.NewVoiceStore(0, IMEI, amrURIPrefix+voiceCacheData.Filename)
			if err != nil {
				log.Info(err)
				fail = true
			}
		}
		if !fail {
			newCmd.AddArg("1")
		} else {
			newCmd.AddArg("2")
		}
		if session != nil {
			if err := session.Send(libnet.Json(newCmd)); err != nil {
				log.Error(err.Error())
			}
		}
	}
	return nil
}

const (
	KC_RAND_KIND_NUM   = 0 // 纯数字
	KC_RAND_KIND_LOWER = 1 // 小写字母
	KC_RAND_KIND_UPPER = 2 // 大写字母
	KC_RAND_KIND_ALL   = 3 // 数字、大小写字母
)

func Krand(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}
