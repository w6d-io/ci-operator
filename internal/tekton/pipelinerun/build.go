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
	"github.com/w6d-io/ci-operator/internal/util"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
)

func (p *PipelineRun) SetBuild(pos int, log logr.Logger) error {
	log = log.WithName("SetBuild").WithValues("action", "pipeline-run")
	log.V(1).Info("set build pipeline run params")

	task := p.Play.Spec.Tasks[pos][ci.Build]
	var flags []string
	if len(task.Variables) != 0 {
		for key, val := range task.Variables {
			flags = append(flags, "--build-arg")
			flags = append(flags, key+"="+val)
		}
	}
	p.Params = append(p.Params, tkn.Param{
		Name: "build_flags",
		Value: tkn.ArrayOrString{
			Type:     tkn.ParamTypeArray,
			ArrayVal: flags,
		},
	})

	if util.IsBuildStage(p.Play) {
		url, err := util.GetDockerImageTag(p.Play)
		if err != nil {
			return err
		}
		p.Params = append(p.Params, tkn.Param{
			Name: "DOCKERFILE",
			Value: tkn.ArrayOrString{
				Type:      tkn.ParamTypeString,
				StringVal: task.Docker.Filepath,
			},
		}, tkn.Param{
			Name: "IMAGE",
			Value: tkn.ArrayOrString{
				Type:      tkn.ParamTypeString,
				StringVal: url.String(),
			},
		}, tkn.Param{
			Name: "CONTEXT",
			Value: tkn.ArrayOrString{
				Type:      tkn.ParamTypeString,
				StringVal: task.Docker.Context,
			},
		})
	}
	return nil
}
