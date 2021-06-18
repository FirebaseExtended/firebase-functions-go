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

package functions

import (
	"context"
	"fmt"

	"github.com/FirebaseExtended/firebase-functions-go/https"
	"github.com/FirebaseExtended/firebase-functions-go/pubsub"
	"github.com/FirebaseExtended/firebase-functions-go/runwith"
)

var Webhook = https.Function{
	RunWith: https.Options{
		AvailableMemoryMB: 256,
	},
	Callback: func(w https.ResponseWriter, r *https.Request) {
		fmt.Fprintf(w, "Hello, world!\n")
	},
}

var PubSubListener = pubsub.Function{
	RunWith: runwith.Options{
		MinInstances: 1,
	},
	EventType: pubsub.V1.Publish,
	Topic:     "topic",
	Callback: func(ctx context.Context, event pubsub.Event) error {
		fmt.Printf("Got event %+v\n", event)
		return nil
	},
}

var NotAFunction = "Non-functions can be safely dumped to emulator.Serve to simplify code gen"

var PubSubListener2 = pubsub.Topic("topic2").OnPublish(func(ctx context.Context, event pubsub.Event) error {
	fmt.Printf("Got event %+v\n", event)
	return nil
})

func NotACloudFunction(x int) {
	fmt.Println(x)
}
