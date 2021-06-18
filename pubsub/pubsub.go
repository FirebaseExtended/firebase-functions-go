// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pubsub

import (
	"errors"
	"fmt"
	"os"

	"github.com/FirebaseExtended/firebase-functions-go/runwith"
	"github.com/FirebaseExtended/firebase-functions-go/support/emulator"
)

type EventType string

var V1 = struct {
	Publish EventType
}{
	Publish: "google.cloud.pubsub.topic.v1.messagePublished",
}

const MessagePublished EventType = "google.pubsub.topic.publish"

type Function struct {
	EventType EventType
	Topic     string
	Region    string
	RunWith   runwith.Options
	Callback  interface{}
}

type Message struct {
	Data       interface{}       `json:"data"`
	Attributes map[string]string `json:"attributes"`
}

type Event struct {
	EventID string  `json:"eventId"`
	Data    Message `json:"data"`
}

func (p Function) AddBackendDescription(symbolName string, b *emulator.Backend) {
	// A builder pattern could ensure Topic is always present...
	if p.Topic == "" {
		panic(fmt.Sprintf("pubsub.Function %s is missing required parameteer Topic", symbolName))
	}
	if p.EventType == "" {
		p.EventType = V1.Publish
	}

	b.CloudFunctions = append(b.CloudFunctions, emulator.FunctionSpec{
		ApiVersion: emulator.GCFv1,
		EntryPoint: fmt.Sprintf("%s.%s", symbolName, "Callback"),
		Id:         symbolName,
		Region:     p.Region,
		Trigger: emulator.EventTrigger{
			EventType: string(p.EventType),
			EventFilters: []emulator.EventFilter{
				{
					Attribute: "resource",
					Value:     fmt.Sprintf("projects/%s/topics/%s", os.Getenv("GCLOUD_PROJECT"), p.Topic),
				},
			},
		},
		MinInstances:      p.RunWith.MinInstances,
		MaxInstances:      p.RunWith.MaxInstances,
		AvailableMemoryMB: p.RunWith.AvailableMemoryMB,
	})
}

func (p Function) Validate() error {
	if p.Topic == "" {
		return errors.New("Pub/Sub functions must define a topic")
	}

	if p.EventType == "" {
		return errors.New("Cloud Functions must have an event type")
	}

	if p.Callback == nil {
		return errors.New("Cloud Functions must have a callback")
	}
	return nil
}

func (p Function) RunWithOptions(options runwith.Options) Function {
	p.RunWith = options
	return p
}

func (p Function) OnPublish(callback interface{}) Function {
	p.EventType = MessagePublished
	p.Callback = callback
	return p
}

func Topic(topicId string) Function {
	return Function{Topic: topicId}
}
