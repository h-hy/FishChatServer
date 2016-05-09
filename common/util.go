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

package common

import (
	"math/rand"
	"time"

	"github.com/oikomi/FishChatServer/base"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/storage/mongo_store"
	"github.com/oikomi/FishChatServer/storage/redis_store"
)

const KeyPrefix string = base.COMM_PREFIX

var DefaultRedisConnectTimeout uint32 = 2000
var DefaultRedisReadTimeout uint32 = 1000
var DefaultRedisWriteTimeout uint32 = 1000

var DefaultRedisOptions redis_store.RedisStoreOptions = redis_store.RedisStoreOptions{
	Network:        "tcp",
	Address:        ":6379",
	ConnectTimeout: time.Duration(DefaultRedisConnectTimeout) * time.Millisecond,
	ReadTimeout:    time.Duration(DefaultRedisReadTimeout) * time.Millisecond,
	WriteTimeout:   time.Duration(DefaultRedisWriteTimeout) * time.Millisecond,
	Database:       1,
	KeyPrefix:      base.COMM_PREFIX,
}

//Just use random to select msg_server
func SelectServer(serverList []string, serverNum int) string {
	return serverList[rand.Intn(serverNum)]
}

func GetSessionFromCID(storeOp interface{}, ID string) (interface{}, error) {
	switch storeOp.(type) {
	case *redis_store.SessionCache:
		// return (*redis_store.SessionCacheData)
		SessionCacheData, err := storeOp.(*redis_store.SessionCache).Get(ID)

		if err != nil {
			log.Warningf("no ID : %s", ID)
			return nil, err
		}
		if SessionCacheData != nil {
			log.Info(SessionCacheData)
		}

		return SessionCacheData, nil
	case *mongo_store.MongoStore:
		// return (*mongo_store.sessionStoreData)
		sessionStoreData, err := storeOp.(*mongo_store.MongoStore).GetSessionFromCid(ID)

		if err != nil {
			log.Warningf("no ID : %s", ID)
			return nil, err
		}
		if sessionStoreData != nil {
			log.Info(sessionStoreData)
		}

		return sessionStoreData, nil

	}

	return nil, NOTFOUNT

	//	session ,err := sessionCache.Get(ID)

	//	if err != nil {
	//		log.Warningf("no ID : %s", ID)
	//		return nil, err
	//	}
	//	if session != nil {
	//		log.Info(session)
	//	}

	//	return session, nil
}

func DelSessionFromCID(storeOp interface{}, ID string) error {
	switch storeOp.(type) {
	case *redis_store.SessionCache:
		err := storeOp.(*redis_store.SessionCache).Delete(ID)

		if err != nil {
			log.Warningf("no ID : %s", ID)
			return err
		}
	}

	//	err := sessionCache.Delete(ID)

	//	if err != nil {
	//		log.Warningf("no ID : %s", ID)
	//		return err
	//	}

	return NOTFOUNT
}

func GetTopicFromTopicName(storeOp interface{}, topicName string) (interface{}, error) {
	switch storeOp.(type) {
	case *redis_store.TopicCache:
		// return (*redis_store.TopicCacheData)
		TopicCacheData, err := storeOp.(*redis_store.TopicCache).Get(topicName)

		if err != nil {
			log.Warningf("no topicName : %s", topicName)
			return nil, err
		}
		if TopicCacheData != nil {
			log.Info(TopicCacheData)
		}

		return TopicCacheData, nil
	case *mongo_store.MongoStore:
		// return (*mongo_store.TopicStoreData)
		TopicStoreData, err := storeOp.(*mongo_store.MongoStore).GetTopicFromCid(topicName)

		if err != nil {
			log.Warningf("no topicName : %s", topicName)
			return nil, err
		}
		if TopicStoreData != nil {
			log.Info(TopicStoreData)
		}

		return TopicStoreData, nil
	}

	return nil, NOTFOUNT

	//	topic ,err := topicCache.Get(topicName)

	//	if err != nil {
	//		log.Warningf("no topicName : %s", topicName)
	//		return nil, err
	//	}
	//	if topic != nil {
	//		log.Info(topic)
	//	}

	//	return topic, nil
}

func StoreData(storeOp interface{}, data interface{}) error {
	switch data.(type) {
	case *redis_store.SessionCacheData:
		err := storeOp.(*redis_store.SessionCache).Set(data.(*redis_store.SessionCacheData))
		//err := (*redis_store.SessionCache)(storeOp).Set(data)
		if err != nil {
			return err
			log.Error("error:", err)
		}
		log.Info("cache sesion success")

		return nil
	case *redis_store.TopicCacheData:
		err := storeOp.(*redis_store.TopicCache).Set(data.(*redis_store.TopicCacheData))
		if err != nil {
			return err
			log.Error("error:", err)
		}
		log.Info("cache topic success")

		return nil
	case *mongo_store.SessionStoreData:
		err := storeOp.(*mongo_store.MongoStore).Set(data)
		if err != nil {
			return err
			log.Error("error:", err)
		}
		log.Info("store session success")

		return nil

	case *mongo_store.TopicStoreData:
		err := storeOp.(*mongo_store.MongoStore).Set(data)
		if err != nil {
			return err
			log.Error("error:", err)
		}
		log.Info("store topic success")

		return nil
	}
	return nil
}

func GetOfflineMsgFromOwnerName(storeOp interface{}, ownerName string) (*redis_store.OfflineMsgCacheData, error) {
	switch storeOp.(type) {
	case *redis_store.OfflineMsgCache:
		o, err := storeOp.(*redis_store.OfflineMsgCache).Get(ownerName)

		if err != nil {
			log.Warningf("no ownerName : %s", ownerName)
			return nil, err
		}
		if o != nil {
			log.Info(o)
		}

		return o, nil
	}

	return nil, NOTFOUNT

	//	o ,err := offlineMsgCache.Get(ownerName)

	//	if err != nil {
	//		log.Warningf("no ownerName : %s", ownerName)
	//		return nil, err
	//	}
	//	if o != nil {
	//		log.Info(o)
	//	}

	//	return o, nil
}
