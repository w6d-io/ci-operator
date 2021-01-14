/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 13/12/2020
*/

package pipeline

import (
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"

	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal/util"
)

func (p *Pipeline) SetPipelineIntTest(play *ci.Play, logger logr.Logger) error {
	log := logger.WithName("SetPipelineIntTest").WithValues("cx-namespace", util.InNamespace(play))

	log.V(1).Info("add task in pipeline")
	wks := getWorkspacePipelineTaskBinding()
	p.Params = append(p.Params, tkn.ParamSpec{
		Name: "init-tests_script",
		Type: tkn.ParamTypeString,
	}, tkn.ParamSpec{
		Name: "init-tests_image",
		Type: tkn.ParamTypeString,
	})
	params := []tkn.Param{
		{
			Name: "script",
			Value: tkn.ArrayOrString{
				Type:      tkn.ParamTypeString,
				StringVal: "$(params.init-tests_script)",
			},
		},
		{
			Name: "IMAGE",
			Value: tkn.ArrayOrString{
				Type:      tkn.ParamTypeString,
				StringVal: "$(params.init-tests_image)",
			},
		},
	}
	task := tkn.PipelineTask{
		Name:       ci.IntegrationTests.String(),
		Resources:  getPipelineTaskResources(false),
		Workspaces: wks,
		Params:     params,
		RunAfter:   p.RunAfter,
		TaskRef: &tkn.TaskRef{
			Kind: tkn.NamespacedTaskKind,
			Name: util.GetCINamespacedName(ci.IntegrationTests.String(), play).Name,
		},
	}
	p.Tasks = append(p.Tasks, task)
	p.RunAfter = append(p.RunAfter, task.Name)
	return nil
}
