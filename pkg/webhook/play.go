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
	"fmt"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

type PlayPayload struct {
	Object     Object    `json:"object,omitempty"`
	ProjectID  int64     `json:"project_id,omitempty"`
	PipelineID int64     `json:"pipeline_id,omitempty"`
	RepoURL    string    `json:"repo_url,omitempty"`
	Commit     ci.Commit `json:"ref,omitempty"`
	Stack      ci.Stack  `json:"stack,omitempty"`
	Status     ci.State  `json:"status,omitempty"`
}

type Object struct {
	NamespacedName types.NamespacedName `json:"namespaced_name,omitempty"`
	Kind           string               `json:"kind,omitempty"`
}

// GetPayLoad returns a filled play payload
func GetPayLoad(play *ci.Play) *PlayPayload {
	nn := getCINamespacedName("pipeline-run", play)
	return &PlayPayload{
		Object: Object{
			Kind:           "pipelinerun",
			NamespacedName: nn,
		},
		ProjectID:  play.Spec.ProjectID,
		PipelineID: play.Spec.PipelineID,
		RepoURL:    play.Spec.RepoURL,
		Commit:     play.Spec.Commit,
		Stack:      play.Spec.Stack,
		Status:     play.Status.State,
	}
}

func getCINamespacedName(prefix string, play *ci.Play) types.NamespacedName {
	return types.NamespacedName{
		Name:      fmt.Sprintf("%s-%v-%v", prefix, play.Spec.ProjectID, play.Spec.PipelineID),
		Namespace: fmt.Sprintf("p6e-cx-%v", play.Spec.ProjectID),
	}
}
