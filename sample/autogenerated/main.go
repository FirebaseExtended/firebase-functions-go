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

// This is a sample of a minimum-complexity boilerplate code
// that must be generated in order to serve the emulator routes
// (GCF and Run) and production routes (Run only) along with
// the admin interface used for backend discovery.
package main

import (
	alias "github.com/FirebaseExtended/firebase-functions-go/sample/lib"
	"github.com/FirebaseExtended/firebase-functions-go/support/emulator"
)

func main() {
	emulator.Serve(map[string]interface{}{
		"Webhook":           alias.Webhook,
		"PubSubListener":    alias.PubSubListener,
		"NotAFunction":      alias.NotAFunction,
		"NotACloudFunction": alias.NotAFunction,
		"PubSubListener2":   alias.PubSubListener2,
	})
}
