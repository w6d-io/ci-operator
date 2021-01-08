/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 15/12/2020
*/

package play

import (
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"

	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal"
	"github.com/w6d-io/ci-operator/internal/tekton/pipelinerun"
	"github.com/w6d-io/ci-operator/internal/util"
)

// SetPipelineRun prepares the PipelineRun and add create method into the run list
func (wf *WFType) SetPipelineRun(play *ci.Play, log logr.Logger) error {
	log = log.WithName("SetPipelineRun").WithValues("cx-namespace", util.InNamespace(play))
	log.Info("Build pipeline run")

	pipelineRun := &pipelinerun.PipelineRun{
		WorkFlowStruct: internal.WorkFlowStruct{
			Scheme: wf.Scheme,
			Play:   play,
		},
	}

	if err := pipelineRun.Parse(log); err != nil {
		return err
	}
	if err := wf.Add(pipelineRun.Create); err != nil {
		return err
	}
	return nil
}
