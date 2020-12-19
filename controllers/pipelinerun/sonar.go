/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 18/12/2020
*/

package pipelinerun

import (
	"github.com/go-logr/logr"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
)

func (p *PipelineRun) SetSonar(pos int, play ci.Play, log logr.Logger) error {
	log = log.WithName("SetBuild").WithValues("action", "pipeline-run")
	log.V(1).Info("set sonar pipeline run params")

	return nil
}
