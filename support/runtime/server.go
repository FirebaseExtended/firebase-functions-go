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

package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"syscall"
	"time"

	"github.com/go-yaml/yaml"
)

type functionData map[string]function

func (f functionData) DescribeBackend(w http.ResponseWriter, r *http.Request) {
	b := Backend{
		SpecVersion:    "v1alpha1",
		RequiredAPIs:   map[string]string{},
		CloudFunctions: make([]FunctionSpec, 0),
		Topics:         make([]PubSubSpec, 0),
		Schedules:      make([]ScheduleSpec, 0),
	}
	for symbol, function := range f {
		function.AddBackendDescription(symbol, &b)
	}
	yaml.NewEncoder(w).Encode(b)
}

func getHandler(f function) func(http.ResponseWriter, *http.Request) {
	// TODO: Should we handle a Callback field that isn't interface{}?
	callback := reflect.ValueOf(f).FieldByName("Callback")
	if callback.Type().Kind() == reflect.Interface {
		callback = callback.Elem()
	}
	if callback.Kind() != reflect.Func {
		panic("CloudFunctions should have a Callback function")
	}
	if callback.Type().NumIn() != 2 {
		panic("CloudFunctions' Callback should take two parameters")
	}

	if httpHandler, ok := callback.Interface().(func(http.ResponseWriter, *http.Request)); ok {
		return httpHandler
	}

	if callback.Type().In(0) != reflect.TypeOf((*context.Context)(nil)).Elem() {
		panic("Event-handling CloudFunctions should take a first parameter of *context.Context")
	}

	if callback.Type().NumOut() != 1 || callback.Type().Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		panic("Event-handling CloudFunctions should return an error")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// Support both a Foo and a Foo* by tracking the arg type and value type separately.
		// If argType is a Foo then valueType is a Foo. If argType is a *Foo then valueType
		// is a Foo.
		argType := callback.Type().In(1)
		valueType := argType
		if valueType.Kind() == reflect.Ptr {
			valueType = valueType.Elem()
		}
		valuePtr := reflect.New(valueType)
		json.NewDecoder(r.Body).Decode(valuePtr.Interface())

		// Now we need to get an actual argument of type argType.
		// valuePtr is type *Foo, so arg will start as Foo and
		// then become *Foo again if argType is a *Foo
		arg := valuePtr.Elem()
		if arg.Type() != argType {
			arg = arg.Addr()
		}

		errVal := callback.Call([]reflect.Value{reflect.ValueOf(r.Context()), arg})[0]
		if errVal.IsNil() {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error handling request: %s", errVal.Interface())
		}
	}
}

// Consider allowing HTTP handler functions to be directly detectable
// and turned into an HttpFunction(-like thing) with default options.
// This would require reimplementing a shim type because we can't have
// a circular reference between the https package and this package.
func Serve(symbols map[string]interface{}) {
	var server http.Server
	var adminServer http.Server

	var port int64 = 8080
	var err error
	if portStr := os.Getenv("PORT"); portStr != "" {
		if port, err = strconv.ParseInt(portStr, 10, 16); err != nil {
			panic("environment variable PORT must be an int")
		}
	}
	server.Addr = fmt.Sprintf("localhost:%d", port)
	fmt.Printf("Serving emulator at http://localhost:%d\n", port)

	if portStr := os.Getenv("ADMIN_PORT"); portStr != "" {
		if adminPort, err := strconv.ParseInt(portStr, 10, 16); err != nil {
			panic("environment varialbe ADMIN_PORT must be an int")
		} else {
			adminServer.Addr = fmt.Sprintf("localhost:%d", adminPort)
			fmt.Printf("Serving emulator admin API at http://localhost:%d\n", adminPort)
		}
	}

	d := functionData{}
	mux := http.NewServeMux()
	server.Handler = mux
	adminMux := http.NewServeMux()
	adminServer.Handler = adminMux

	for symbol, value := range symbols {
		if asFunc, ok := value.(function); ok {
			d[symbol] = asFunc
		}
	}

	for symbol, function := range d {
		fmt.Printf("Serving function at http://localhost:%d/%s\n", port, symbol)
		mux.HandleFunc(fmt.Sprintf("/%s", symbol), getHandler(function))
	}

	shouldShutDown := make(chan os.Signal, 1)
	signal.Notify(shouldShutDown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	adminMux.HandleFunc("/backend.yaml", d.DescribeBackend)
	adminMux.HandleFunc("/quitquitquit", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK\n")
		shouldShutDown <- syscall.SIGINT
	})

	didShutDown := make(chan struct{}, 2)
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Println("Emulator exited with error", err)
		}
		didShutDown <- struct{}{}
	}()
	go func() {
		if adminServer.Addr != "" {
			err := adminServer.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				fmt.Println("Emulator admin API exited with error", err)
			}
		}
		didShutDown <- struct{}{}
	}()

	<-shouldShutDown
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	server.Shutdown(ctx)
	adminServer.Shutdown(ctx)
	<-didShutDown
	<-didShutDown
}

type function interface {
	AddBackendDescription(symbolName string, b *Backend)
}
