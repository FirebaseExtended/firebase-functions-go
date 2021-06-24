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
	"os"
)

type EventFilters []EventFilter

func (f EventFilters) MarshalYAML() (interface{}, error) {
	m := make(map[string]interface{}, len(f))
	for _, filter := range f {
		m[filter.Attribute] = filter.Value
	}
	return m, nil
}

type EventFilter struct {
	Attribute string `yaml:"attribute"`
	Value     string `yaml:"value"`
}

type EventTrigger struct {
	EventType           string       `yaml:"eventType,omitempty"`
	EventFilters        EventFilters `yaml:"eventFilters,omitempty"`
	ServiceAccountEmail string       `yaml:"serviceAccountEmail,omitempty"`
}

type ApiVersion int

const GCFv2 ApiVersion = 2
const GCFv1 ApiVersion = 1

type FunctionSpec struct {
	ApiVersion ApiVersion `yaml:"apiVersion"`
	EntryPoint string     `yaml:"entryPoint"`
	Id         string     `yaml:"id"`
	Region     string     `yaml:"region,omitempty"`
	Project    string     `yaml:"project,omitempty"`
	// NOTE: In the current schema this is a union between
	// an HTTP and an EventTrigger. Since HTTP triggers have
	// no options in GCFv2 we can just use an empty EventTrigger
	// for now.
	Trigger           EventTrigger `yaml:"trigger"`
	MinInstances      int          `yaml:"minInstances,omitempty"`
	MaxInstances      int          `yaml:"maxInstances,omitempty"`
	AvailableMemoryMB int          `yaml:"availableMemoryMb,omitempty"`
}

type TargetService struct {
	Id      string `yaml:"id"`
	Region  string `yaml:"region,omitempty"`
	Project string `yaml:"project,omitempty"`
}

type PubSubSpec struct {
	Id            string        `yaml:"id"`
	Project       string        `yaml:"project,omitempty"`
	TargetService TargetService `yaml:"targetService"`
}

type ScheduleRetryConfig struct {
	RetryCount int `yaml:"retryCount,omitempty"`
}

type Transport string

const PubSubTransport Transport = "pubsub"
const HttpsTransport Transport = "https"

type ScheduleSpec struct {
	Id            string              `yaml:"id"`
	Project       string              `yaml:"project"`
	Schedule      string              `yaml:"schedule"`
	TimeZone      string              `yaml:"timeZone,omitempty"`
	RetryConfig   ScheduleRetryConfig `yaml:"retryConfig"`
	Transport     Transport           `yaml:"transport"`
	TargetService TargetService       `yaml:"targetService"`
}

type Backend struct {
	SpecVersion    string            `yaml:"specVersion"`
	RequiredAPIs   map[string]string `yaml:"requiredAPIs,omitempty"`
	CloudFunctions []FunctionSpec    `yaml:"cloudFunctions"`
	Topics         []PubSubSpec      `yaml:"topics,omitempty"`
	Schedules      []ScheduleSpec    `yaml:"schedules,omitempty"`
}

func ProjectOrDefault(project string) string {
	if project != "" {
		return project
	}
	return os.Getenv("GCLOUD_PROJECT")
}
