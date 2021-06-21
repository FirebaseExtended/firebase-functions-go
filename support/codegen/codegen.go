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

package main

import (
	"errors"
	"fmt"
	"go/types"
	"log"
	"os"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) != 1 {
		fmt.Println("Usage: go run codegen.go <FunctionPackage>")
		os.Exit(1)
	}

	pkg := argsWithoutProg[0]
	triggers, err := extractTriggers(pkg)
	if err != nil {
		log.Fatalln(err)
	}

	code, err := generateEntrypoint(pkg, triggers)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(code)
}

func extractTriggers(pkg string) ([]string, error) {
	exports, err := loadExports(pkg)
	if err != nil {
		return nil, err
	}

	var triggers []string
	for _, exp := range exports {
		triggers = append(triggers, exp.Name())
	}

	return triggers, nil
}

func loadExports(pkg string) ([]types.Object, error) {
	cfg := &packages.Config{Mode: packages.LoadTypes}
	pkgs, err := packages.Load(cfg, pkg)
	if err != nil {
		return nil, err
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, fmt.Errorf("one or mor errors")
	}

	var results []types.Object
	scope := pkgs[0].Types.Scope()
	for _, symbol := range scope.Names() {
		typeInfo := scope.Lookup(symbol)
		if typeInfo.Exported() {
			results = append(results, typeInfo)
		}
	}

	if results == nil {
		return nil, errors.New("no exported member")
	}

	return results, nil
}

const mainTemplate = `
package main

import (
	alias "{{ .Pkg }}"
	"github.com/FirebaseExtended/firebase-functions-go/support/emulator"
)

func main() {
	emulator.Serve(map[string]interface{}{
	{{- range .Triggers }}
		"{{ . }}": alias.{{ . }},
	{{- end }}
	})
}
`

func generateEntrypoint(pkg string, triggers []string) (string, error) {
	tmpl, err := template.New("main").Parse(mainTemplate)
	if err != nil {
		return "", err
	}

	b := new(strings.Builder)
	info := struct {
		Pkg      string
		Triggers []string
	}{
		Pkg:      pkg,
		Triggers: triggers,
	}
	if err := tmpl.Execute(b, info); err != nil {
		return "", err
	}

	return b.String(), nil
}
