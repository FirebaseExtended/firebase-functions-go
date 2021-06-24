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

package https

import (
	"fmt"
	"net/http"

	"github.com/FirebaseExtended/firebase-functions-go/support/runtime"
)

type Request = http.Request
type ResponseWriter = http.ResponseWriter

type Options struct {
	MinInstances      int
	MaxInstances      int
	AvailableMemoryMB int
}

type Function struct {
	Callback func(ResponseWriter, *Request)
	RunWith  Options
}

func (h Function) AddBackendDescription(symbolName string, b *runtime.Backend) {
	// Runtime isn't specified from within the API?
	b.CloudFunctions = append(b.CloudFunctions, runtime.FunctionSpec{
		ApiVersion:        runtime.GCFv1,
		Id:                symbolName,
		EntryPoint:        fmt.Sprintf("%s.%s", symbolName, "Callback"),
		MinInstances:      h.RunWith.MinInstances,
		MaxInstances:      h.RunWith.MaxInstances,
		AvailableMemoryMB: h.RunWith.AvailableMemoryMB,
	})
}

func (h Function) Valdiate() error {
	return nil
}

func (h Function) RunWithOptions(options Options) Function {
	h.RunWith = options
	return h
}
