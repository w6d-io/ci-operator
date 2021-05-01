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
	"fmt"

	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/util"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
)

// SetDeploy builds the params adn add the minio volume for the tekton pipelineRun resource
// param pos : use to get task part from Play to get variables
func (p *PipelineRun) SetDeploy(pos int, log logr.Logger) error {
	log = log.WithName("SetDeploy").WithValues("action", "pipeline-run")
	log.V(1).Info("set deploy pipeline run params")

	task := p.Play.Spec.Tasks[pos][ci.Deploy]
	var flags []string
	if len(task.Variables) != 0 {
		for key, val := range task.Variables {
			flags = append(flags, "--set")
			flags = append(flags, key+"="+val)
		}
	}
	p.Params = append(p.Params, tkn.Param{
		Name: "deploy_flags",
		Value: tkn.ArrayOrString{
			Type:     tkn.ParamTypeArray,
			ArrayVal: flags,
		},
	}, tkn.Param{
		Name: "deploy_s3valuepath",
		Value: tkn.ArrayOrString{
			Type: tkn.ParamTypeString,
			StringVal: fmt.Sprintf("%v/%v/%v/values.yaml",
				config.GetMinio().Bucket,
				p.Play.Spec.ProjectID,
				p.Play.Spec.PipelineID),
		},
	}, tkn.Param{
		Name: "deploy_values",
		Value: tkn.ArrayOrString{
			Type:      tkn.ParamTypeString,
			StringVal: fmt.Sprintf("%s/values.yaml", config.GetWorkspacePath("values")),
		},
	}, tkn.Param{
		Name: "deploy_namespace",
		Value: tkn.ArrayOrString{
			Type: tkn.ParamTypeString,
			// TODO put the prefix in config
			StringVal: p.GetNamespace(task),
		},
	}, tkn.Param{
		Name: "deploy_release_name",
		Value: tkn.ArrayOrString{
			Type:      tkn.ParamTypeString,
			StringVal: p.GetNamespace(task),
		},
	})

	return nil
}

func (p *PipelineRun) GetNamespace(task ci.Task) string {
	if task.Namespace != "" {
		return task.Namespace
	}
	return util.GetDeployNamespacedName("cx", p.Play).Namespace
}
