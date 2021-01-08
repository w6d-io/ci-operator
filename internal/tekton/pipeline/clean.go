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
	"github.com/go-logr/logr"
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal/util"
)

func (p *Pipeline) SetPipelineClean(play *ci.Play, logger logr.Logger) error {
	log := logger.WithName("SetPipelineClean").WithValues("cx-namespace", util.InNamespace(play))

	log.V(1).Info("add task in pipeline")
	wks := getWorkspacePipelineTaskBinding()
	p.Params = append(p.Params, tkn.ParamSpec{
		Name: "clean_namespace",
		Type: tkn.ParamTypeString,
	}, tkn.ParamSpec{
		Name: "clean_release_name",
		Type: tkn.ParamTypeString,
	})
	params := []tkn.Param{
		{
			Name: "namespace",
			Value: tkn.ArrayOrString{
				Type:      tkn.ParamTypeString,
				StringVal: "$(params.clean_namespace)",
			},
		},
		{
			Name: "release_name",
			Value: tkn.ArrayOrString{
				Type:      tkn.ParamTypeString,
				StringVal: "$(params.clean_release_name)",
			},
		},
	}
	task := tkn.PipelineTask{
		Name:       ci.Clean.String(),
		Workspaces: wks,
		Params:     params,
		RunAfter:   p.RunAfter,
		TaskRef: &tkn.TaskRef{
			Kind: tkn.NamespacedTaskKind,
			Name: util.GetCINamespacedName(ci.Clean.String(), play).Name,
		},
	}
	p.Tasks = append(p.Tasks, task)
	p.RunAfter = append(p.RunAfter, task.Name)
	return nil
}
