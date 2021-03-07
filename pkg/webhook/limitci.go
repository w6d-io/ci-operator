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
Created on 04/03/2021
*/
package webhook

import (
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
)

type Message struct {
	ProjectID int64        `json:"projectid"`
	EventType string       `json:"message_type"`
	Data      LimitPayload `json:"data"`
}

type LimitPayload struct {
	Object     Object `json:"object,omitempty"`
	Concurrent int64  `json:"concurrent,omitempty"`
	Message    string `json:"message,omitempty"`
}

// GetLimitPayload return a filled limit ci payload
func GetLimitPayload(play *ci.Play, limit ci.LimitCi, message string) *Message {
	nn := getCINamespacedName("pipeline-run", play)
	return &Message{
		ProjectID: play.Spec.ProjectID,
		EventType: "ci-operator",
		Data: LimitPayload{
			Object: Object{
				Kind:           "pipelinerun",
				NamespacedName: nn,
			},
			Concurrent: limit.Spec.Concurrent,
			Message:    message,
		}}
}
