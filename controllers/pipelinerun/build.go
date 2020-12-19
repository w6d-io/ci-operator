/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 17/12/2020
*/

package pipelinerun

import (
	"github.com/go-logr/logr"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
)

func (p *PipelineRun) SetBuild(pos int, play ci.Play, log logr.Logger) error {
	log = log.WithName("SetBuild").WithValues("action", "pipeline-run")
	log.V(1).Info("set build pipeline run params")

	task := play.Spec.Tasks[pos][ci.Build]
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

	if p.IsBuildStage(play) {
		url, err := p.GetDockerImage(play)
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
