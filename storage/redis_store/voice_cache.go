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

package redis_store

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

type VoiceCache struct {
	RS      *RedisStore
	rwMutex sync.Mutex
}

func NewVoiceCache(RS *RedisStore) *VoiceCache {
	return &VoiceCache{
		RS: RS,
	}
}

type VoiceCacheData struct {
	Type         string
	Id           int
	Uri          string
	PathFilename string
	Format       string
	Size         int
	NowGroup     int
	NowSize      int
}

func (self *VoiceCache) NewVoiceCacheData(Type string, id int, uri, pathFilename, format string, size int) *VoiceCacheData {

	cacheData := &VoiceCacheData{
		Type:         Type,
		Id:           id,
		Uri:          uri,
		PathFilename: pathFilename,
		Format:       format,
		Size:         size,
		NowGroup:     0,
		NowSize:      0,
	}
	return cacheData
}

// Get the session from the store.
func (self *VoiceCache) Get(Type, voiceId string) (*VoiceCacheData, error) {
	self.rwMutex.Lock()
	defer self.rwMutex.Unlock()
	key := "VOICE" + Type + ":" + voiceId
	log.Println(key)
	b, err := redis.Bytes(self.RS.conn.Do("GET", key))
	if err != nil {
		return nil, err
	}
	var sess VoiceCacheData
	err = json.Unmarshal(b, &sess)
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

// Save the session into the store.
func (self *VoiceCache) Set(sess *VoiceCacheData) error {
	self.rwMutex.Lock()
	defer self.rwMutex.Unlock()
	b, err := json.Marshal(sess)
	if err != nil {
		return err
	}
	key := "VOICE" + sess.Type + ":" + strconv.Itoa(sess.Id)
	fmt.Print(key)
	ttl := 600 * time.Second
	_, err = self.RS.conn.Do("SETEX", key, int(ttl.Seconds()), b)
	if err != nil {
		return err
	}
	return nil
}

// Delete the session from the store.
func (self *VoiceCache) Delete(voiceId string) error {
	self.rwMutex.Lock()
	defer self.rwMutex.Unlock()
	key := "VOICE:" + voiceId
	_, err := self.RS.conn.Do("DEL", key)
	if err != nil {
		return err
	}
	return nil
}
