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
Created on 01/03/2021
*/
package pipelinerun

import (
    "github.com/go-logr/logr"
    tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
    ci "github.com/w6d-io/ci-operator/api/v1alpha1"
    "strings"
)

func (p *PipelineRun) SetE2ETest(pos int, logger logr.Logger) error {
    log := logger.WithValues("SetE2ETest").WithValues("action", "pipeline-run")
    log.V(1).Info("set e2e test pipeline run params")

    task := p.Play.Spec.Tasks[pos][ci.E2ETests]
    p.Params = append(p.Params, tkn.Param{
        Name: "e2e-tests_script",
        Value: tkn.ArrayOrString{
            Type:      tkn.ParamTypeString,
            StringVal: strings.Join(task.Script, "\n"),
        },
    }, tkn.Param{
        Name:  "e2e-tests_image",
        Value: tkn.ArrayOrString{
            Type:      tkn.ParamTypeString,
            StringVal: task.Image,
        },
    })
    return nil
}