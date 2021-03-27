/*
Copyright 2020 WILDCARD

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
Created on 15/12/2020
*/

package play

import (
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"

	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal"
	"github.com/w6d-io/ci-operator/internal/tekton/pipelinerun"
)

// SetPipelineRun prepares the PipelineRun and add create method into the run list
func (wf *WFType) SetPipelineRun(play *ci.Play, log logr.Logger) error {
	log = log.WithName("SetPipelineRun")
	log.Info("Build pipeline run")

	pipelineRun := &pipelinerun.PipelineRun{
		WorkFlowStruct: internal.WorkFlowStruct{
			Scheme: wf.Scheme,
			Play:   play,
		},
	}

	if err := pipelineRun.Parse(log); err != nil {
		return err
	}
	if err := wf.Add(pipelineRun.Create); err != nil {
		return err
	}
	return nil
}
