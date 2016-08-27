package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/oikomi/FishChatServer/log"
)

type Voice struct {
	Type         string
	Id           int
	Uri          string
	Filename     string
	PathFilename string
	Format       string
	Size         int
}

func (this *Voice) Cache() {
	body, err := json.Marshal(this)
	log.Info(err)
	if err == nil {
		redisCache.Put("VOICEDOWN:"+strconv.Itoa(this.Id), body, 30*time.Minute)
	}
}

//func getVoice(Type,Id int) (Voice, error) {
//	voice := redisCache.Get("VOICE"+Type+":" + string(Id))
//	if voice != nil {
//		voiceString := GetString(voice)
//		var voiceObj Voice
//		err := json.Unmarshal([]byte(voiceString), &voiceObj)
//		return voiceObj, err
//	}
//	return Voice{}, nil
//}

func GetString(v interface{}) string {
	switch result := v.(type) {
	case string:
		return result
	case []byte:
		return string(result)
	default:
		if v != nil {
			return fmt.Sprintf("%v", result)
		}
	}
	return ""
}
