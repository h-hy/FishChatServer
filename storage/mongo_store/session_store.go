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

type SessionStoreData struct {
	ClientID   string   `bson:"ClientID"`
	ClientPwd  string   `bson:"ClientPwd"`
	ClientName string   `bson:"ClientName"`
	ClientType string   `bson:"ClientType"`
	TopicList  []string `bson:"TopicList"`
}

func NewSessionStoreData(clientID string, clientPwd string, clientType string) *SessionStoreData {
	return &SessionStoreData{
		ClientID:   clientID,
		ClientPwd:  clientPwd,
		ClientType: clientType,
		TopicList:  make([]string, 0),
	}
}

func (self *SessionStoreData) AddTopic(t string) {
	self.TopicList = append(self.TopicList, t)
}

func (self *SessionStoreData) GetTopics() []string {
	return self.TopicList
}

func (self *SessionStoreData) TopicExist(t string) bool {
	for _, name := range self.TopicList {
		if name == t {
			return true
		}
	}
	return false
}

func (self *SessionStoreData) RemoveTopic(t string) {
	for i, name := range self.TopicList {
		if name == t {
			self.TopicList = append(self.TopicList[:i], self.TopicList[i+1:]...)
			break
		}
	}
}
