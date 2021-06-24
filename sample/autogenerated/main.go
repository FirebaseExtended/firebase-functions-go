
package main

import (
	alias "github.com/FirebaseExtended/firebase-functions-go/sample/lib"
	"github.com/FirebaseExtended/firebase-functions-go/support/runtime"
)

func main() {
	runtime.Serve(map[string]interface{}{
		"NotACloudFunction": alias.NotACloudFunction,
		"NotAFunction": alias.NotAFunction,
		"PubSubListener": alias.PubSubListener,
		"PubSubListener2": alias.PubSubListener2,
		"Webhook": alias.Webhook,
	})
}

