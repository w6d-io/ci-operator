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

// SetPipelineBuild adds build tasks elements in pipeline
func (p *Pipeline) SetPipelineBuild(logger logr.Logger) error {

	log := logger.WithName("SetPipelineBuild")
	tk := p.Play.Spec.Tasks[p.Pos][ci.Build]
	var flags []string
	if len(tk.Variables) != 0 {
		for key, val := range tk.Variables {
			flags = append(flags, "--build-arg")
			flags = append(flags, key+"="+val)
		}
	}
	log.V(1).Info("add task in pipeline")
	wks := getWorkspacePipelineTaskBinding()
	p.Params = append(p.Params, tkn.ParamSpec{
		Name: "build_flags",
		Type: tkn.ParamTypeArray,
	}, tkn.ParamSpec{
		Name: "build_s3DockerfilePath",
		Type: tkn.ParamTypeString,
	})
	params := []tkn.Param{
		{
			Name: "flags",
			Value: tkn.ArrayOrString{
				Type:     tkn.ParamTypeArray,
				ArrayVal: flags,
			},
		},
		{
			Name: "s3DockerfilePath",
			Value: tkn.ArrayOrString{
				Type:      tkn.ParamTypeString,
				StringVal: "$(params.build_s3DockerfilePath)",
			},
		},
	}
	if util.IsBuildStage(p.Play) {
		log.V(1).Info("add parameter")
		params = append(params,
			getParamString("DOCKERFILE"),
			getParamString("IMAGE"),
			getParamString("CONTEXT"),
		)
		p.Params = append(p.Params, tkn.ParamSpec{
			Name: "DOCKERFILE",
			Type: tkn.ParamTypeString,
		}, tkn.ParamSpec{
			Name: "IMAGE",
			Type: tkn.ParamTypeString,
		}, tkn.ParamSpec{
			Name: "CONTEXT",
			Type: tkn.ParamTypeString,
		})
	}
	task := tkn.PipelineTask{
		Name:       ci.Build.String(),
		Resources:  getPipelineTaskResources(util.IsBuildStage(p.Play)),
		Workspaces: wks,
		Params:     params,
		RunAfter:   p.RunAfter,
		TaskRef: &tkn.TaskRef{
			Kind: tkn.NamespacedTaskKind,
			Name: util.GetCINamespacedName(ci.Build.String(), p.Play).Name,
		},
	}
	p.Tasks = append(p.Tasks, task)
	p.RunAfter = append(p.RunAfter, task.Name)
	return nil
}
