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
Created on 17/12/2020
*/

package pipelinerun

import (
	"github.com/go-logr/logr"
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
)

func (p *PipelineRun) SetGeneric(pos int, taskType ci.TaskType, log logr.Logger) error {
	log = log.WithName("SetBuild").
		WithValues("action", "pipeline-run", "task-type", taskType)
	log.V(1).Info("set build pipeline run params")

	task := p.Play.Spec.Tasks[pos][taskType]
	var flags []string
	if len(task.Arguments) != 0 {
		for _, val := range task.Arguments {
			flags = append(flags, val)
		}
	}
	p.Params = append(p.Params, tkn.Param{
		Name: string(taskType) + "_flags",
		Value: tkn.ArrayOrString{
			Type:     tkn.ParamTypeArray,
			ArrayVal: flags,
		},
	})
	for _, gp := range p.GenericParams[string(taskType)] {
		p.Params = append(p.Params, tkn.Param{
			Name: string(taskType) + "_" + gp.Name,
			Value: tkn.ArrayOrString{
				Type:      tkn.ParamTypeString,
				StringVal: p.Play.Get(gp.Value),
			},
		})
	}
	return nil
}
