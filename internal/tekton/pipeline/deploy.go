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
Created on 13/12/2020
*/

package pipeline

import (
	"github.com/go-logr/logr"
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal/util"
)

func (p *Pipeline) SetPipelineDeploy(logger logr.Logger) error {
	log := logger.WithName("SetPipelineDeploy").WithValues("cx-namespace", util.InNamespace(p.Play))

	log.V(1).Info("add task in pipeline")
	wks := getWorkspacePipelineTaskBinding()
	p.Params = append(p.Params, tkn.ParamSpec{
		Name: "deploy_flags",
		Type: tkn.ParamTypeArray,
	}, tkn.ParamSpec{
		Name: "deploy_s3valuepath",
		Type: tkn.ParamTypeString,
	}, tkn.ParamSpec{
		Name: "deploy_values",
		Type: tkn.ParamTypeString,
	}, tkn.ParamSpec{
		Name: "deploy_namespace",
		Type: tkn.ParamTypeString,
	}, tkn.ParamSpec{
		Name: "deploy_release_name",
		Type: tkn.ParamTypeString,
	})
	params := []tkn.Param{
		{
			Name: "flags",
			Value: tkn.ArrayOrString{
				Type:     tkn.ParamTypeArray,
				ArrayVal: []string{"$(params.deploy_flags)"},
			},
		},
		{
			Name: "s3valuepath",
			Value: tkn.ArrayOrString{
				Type:      tkn.ParamTypeString,
				StringVal: "$(params.deploy_s3valuepath)",
			},
		},
		{
			Name: "values",
			Value: tkn.ArrayOrString{
				Type:      tkn.ParamTypeString,
				StringVal: "$(params.deploy_values)",
			},
		},
		{
			Name: "namespace",
			Value: tkn.ArrayOrString{
				Type:      tkn.ParamTypeString,
				StringVal: "$(params.deploy_namespace)",
			},
		},
		{
			Name: "release_name",
			Value: tkn.ArrayOrString{
				Type:      tkn.ParamTypeString,
				StringVal: "$(params.deploy_release_name)",
			},
		},
	}

	task := tkn.PipelineTask{
		Name:       ci.Deploy.String(),
		Workspaces: wks,
		Params:     params,
		RunAfter:   p.RunAfter,
		TaskRef: &tkn.TaskRef{
			Kind: tkn.NamespacedTaskKind,
			Name: util.GetCINamespacedName(ci.Deploy.String(), p.Play).Name,
		},
	}
	p.Tasks = append(p.Tasks, task)
	p.RunAfter = append(p.RunAfter, task.Name)
	return nil
}
