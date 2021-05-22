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
Created on 12/12/2020
*/

package pipeline

import (
	"github.com/go-logr/logr"
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal/util"
)

// SetPipelineGeneric adds build tasks elements in pipeline
func (p *Pipeline) SetPipelineGeneric(taskType ci.TaskType, logger logr.Logger) error {

	log := logger.WithName("SetPipelineGeneric").WithValues("task", taskType)
	tk := p.Play.Spec.Tasks[p.Pos][taskType]
	var flags []string
	if len(tk.Arguments) != 0 {
		for _, val := range tk.Arguments {
			flags = append(flags, val)
		}
	}
	log.V(1).Info("get workspace")
	wks := getWorkspacePipelineTaskBinding()
	log.V(1).Info("build params")
	params := p.BuildParams(taskType)
	log.V(1).Info("add task in pipeline")
	task := tkn.PipelineTask{
		Name:       taskType.String(),
		Resources:  getPipelineTaskResources(false),
		Workspaces: wks,
		Params:     params,
		RunAfter:   p.RunAfter,
		TaskRef: &tkn.TaskRef{
			Kind: tkn.NamespacedTaskKind,
			Name: util.GetCINamespacedName(taskType.String(), p.Play).Name,
		},
	}
	p.Tasks = append(p.Tasks, task)
	p.RunAfter = append(p.RunAfter, task.Name)
	return nil
}

func (p *Pipeline) BuildParams(taskType ci.TaskType) (params []tkn.Param) {
	p.Params = append(p.Params, tkn.ParamSpec{
		Name: string(taskType) + "_flags",
		Type: tkn.ParamTypeArray,
	})
	params = []tkn.Param{
		{
			Name: "flags",
			Value: tkn.ArrayOrString{
				Type:     tkn.ParamTypeArray,
				ArrayVal: []string{"$(params." + string(taskType) + "_flags)"},
			},
		},
	}
	for _, gp := range p.GenericParams[string(taskType)] {
		p.Params = append(p.Params, tkn.ParamSpec{
			Name: string(taskType) + "_" + gp.Name,
			Type: tkn.ParamTypeString,
		})

		params = append(params, tkn.Param{
			Name: gp.Name,
			Value: tkn.ArrayOrString{
				Type:      gp.Type,
				StringVal: "$(params." + string(taskType) + "_" + gp.Name + ")",
			},
		})
	}
	return
}
