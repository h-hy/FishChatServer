//
// Copyright 2014-2015 Hong Miao (miaohong@miaohong.org). All Rights Reserved.
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

package mongo_store

import (
	"sync"
	"time"

	"github.com/oikomi/FishChatServer/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoStoreOptions struct {
}

type MongoStore struct {
	opts    *MongoStoreOptions
	session *mgo.Session

	rwMutex sync.Mutex
}

func NewMongoStore(ip string, port string, user string, password string) *MongoStore {
	var url string
	if user == "" && password == "" {
		url = ip + port
	} else {
		url = user + ":" + password + "@" + ip + port
	}

	log.Info("connect to mongo : ", url)
	maxWait := time.Duration(5 * time.Second)
	session, err := mgo.DialWithTimeout(url, maxWait)
	session.SetMode(mgo.Monotonic, true)
	if err != nil {
		panic(err)
	}
	return &MongoStore{
		session: session,
	}
}

func (self *MongoStore) Init() {
	//self.session.DB("im").C("client_info")

}

func (self *MongoStore) Set(data interface{}) error {
	log.Info("MongoStore Update")
	var err error
	self.rwMutex.Lock()
	defer self.rwMutex.Unlock()

	switch data.(type) {
	case *SessionStoreData:
		op := self.session.DB(DATA_BASE_NAME).C(CLIENT_INFO_COLLECTION)
		cid := data.(*SessionStoreData).ClientID
		log.Info("cid : ", cid)
		_, err = op.Upsert(bson.M{"ClientID": cid}, data.(*SessionStoreData))
		if err != nil {
			log.Error(err.Error())
			return err
		}

	case *TopicStoreData:
		op := self.session.DB(DATA_BASE_NAME).C(TOPIC_INFO_COLLECTION)
		topicName := data.(*TopicStoreData).TopicName
		log.Info("topicName : ", topicName)
		_, err = op.Upsert(bson.M{"TopicName": topicName}, data.(*TopicStoreData))
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	return err
}

func (self *MongoStore) GetTopicFromCid(cid string) (*TopicStoreData, error) {
	log.Info("MongoStore GetTopicFromCid")
	var err error
	self.rwMutex.Lock()
	defer self.rwMutex.Unlock()

	op := self.session.DB(DATA_BASE_NAME).C(TOPIC_INFO_COLLECTION)

	var result *TopicStoreData

	err = op.Find(bson.M{"TopicName": cid}).One(result)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return result, nil
}

func (self *MongoStore) DeleteTopic(cid string) error {
	log.Info("MongoStore DeleteTopic")

	self.rwMutex.Lock()
	defer self.rwMutex.Unlock()

	op := self.session.DB(DATA_BASE_NAME).C(TOPIC_INFO_COLLECTION)

	return op.Remove(bson.M{"TopicName": cid})
}

func (self *MongoStore) GetSessionFromCid(cid string) (*SessionStoreData, error) {
	log.Info("MongoStore GetSessionFromCid")
	var err error
	self.rwMutex.Lock()
	defer self.rwMutex.Unlock()

	op := self.session.DB(DATA_BASE_NAME).C(CLIENT_INFO_COLLECTION)

	var result *SessionStoreData

	err = op.Find(bson.M{"ClientID": cid}).One(result)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return result, nil
}

func (self *MongoStore) DeleteSession(cid string) error {
	log.Info("MongoStore DeleteSession")

	self.rwMutex.Lock()
	defer self.rwMutex.Unlock()

	op := self.session.DB(DATA_BASE_NAME).C(CLIENT_INFO_COLLECTION)

	return op.Remove(bson.M{"ClientID": cid})
}

func (self *MongoStore) Close() {
	self.session.Close()
}
