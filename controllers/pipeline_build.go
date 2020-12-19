/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 12/12/2020
*/

package controllers

import (
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
)

// SetBuild adds build tasks elements in pipeline
func (p *Pipeline) SetPipelineBuild(play ci.Play, r *PlayReconciler) error {

	log := r.Log.WithName("SetPipelineBuild").WithValues("cx-namespace", InNamespace(play))

	log.V(1).Info("add task in pipeline")
	wks := getWorkspacePipelineTaskBinding()
	p.Params = append(p.Params, tkn.ParamSpec{
		Name: "build_flags",
		Type: tkn.ParamTypeArray,
	})
	params := []tkn.Param{
		{
			Name: "flags",
			Value: tkn.ArrayOrString{
				Type:     tkn.ParamTypeArray,
				ArrayVal: []string{"$(params.build_flags[*])"},
			},
		},
	}
	if IsBuildStage(play) {
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
		Resources:  getPipelineTaskResources(IsBuildStage(play)),
		Workspaces: wks,
		Params:     params,
		RunAfter:   p.RunAfter,
		TaskRef: &tkn.TaskRef{
			Kind: tkn.NamespacedTaskKind,
			Name: CxCINamespacedName(ci.Build.String(), play).Name,
		},
	}
	p.Tasks = append(p.Tasks, task)
	p.RunAfter = append(p.RunAfter, task.Name)
	return nil
}
