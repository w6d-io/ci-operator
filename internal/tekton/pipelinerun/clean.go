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
	"github.com/w6d-io/ci-operator/internal/util"
)

func (p *PipelineRun) SetClean(_ int, log logr.Logger) error {
	log = log.WithName("SetClean").WithValues("action", "pipeline-run")
	log.V(1).Info("set clean pipeline run params")
	p.Params = append(p.Params, tkn.Param{
		Name: "clean_namespace",
		Value: tkn.ArrayOrString{
			Type:      tkn.ParamTypeString,
			StringVal: util.GetDeployNamespacedName("cx", p.Play).Namespace,
		},
	}, tkn.Param{
		Name: "clean_release_name",
		Value: tkn.ArrayOrString{
			Type:      tkn.ParamTypeString,
			StringVal: util.GetDeployNamespacedName("cx", p.Play).Namespace,
		},
	})
	return nil
}
