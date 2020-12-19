/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 14/12/2020
*/

package controllers

import (
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
)

func (p *Pipeline) SetPipelineUnitTest(play ci.Play, r *PlayReconciler) error {
	log := r.Log.WithName("SetPipelineUnitTest").WithValues("cx-namespace", InNamespace(play))

	log.V(1).Info("add task in pipeline")
	wks := getWorkspacePipelineTaskBinding()
	p.Params = append(p.Params, tkn.ParamSpec{
		Name: "unit-tests_script",
		Type: tkn.ParamTypeString,
	})
	params := []tkn.Param{
		{
			Name: "script",
			Value: tkn.ArrayOrString{
				Type:      tkn.ParamTypeString,
				StringVal: "$(params.unit-tests_script)",
			},
		},
	}
	task := tkn.PipelineTask{
		Name:       ci.UnitTests.String(),
		Resources:  getPipelineTaskResources(IsBuildStage(play)),
		Workspaces: wks,
		Params:     params,
		RunAfter:   p.RunAfter,
		TaskRef: &tkn.TaskRef{
			Kind: tkn.NamespacedTaskKind,
			Name: CxCINamespacedName(ci.UnitTests.String(), play).Name,
		},
	}
	p.Tasks = append(p.Tasks, task)
	p.RunAfter = append(p.RunAfter, task.Name)
	return nil

}
