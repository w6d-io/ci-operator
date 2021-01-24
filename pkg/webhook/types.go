/*
Copyright 2020 WILDCARD SA.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
Created on 20/01/2021
*/
package webhook

import (
	"net/url"

	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
)

type Payload interface {
	// Send sends the event to a webhook address
	Send(string) error

	// DoSend loop on webhook address and call Send
	DoSend([]Webhook) error

	// SetStatus record the status in the Payload
	SetStatus(ci.State)

	// SetStatus record the status in the Payload
	SetObjectName(string)
}

var payload Payload

type Webhook struct {
	Name   string `json:"name" yaml:"name"`
	URLRaw string `json:"url" yaml:"url"`
	URL    *url.URL
}
