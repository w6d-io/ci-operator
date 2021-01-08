/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 08/12/2020
*/

package play

import (
	"github.com/go-logr/logr"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal/tekton/pipeline"
	"github.com/w6d-io/ci-operator/internal/util"
)

func (wf *WFType) SetPipeline(p *ci.Play, logger logr.Logger) error {
	log := logger.WithName("SetPipeline").WithValues("cx-namespace", util.InNamespace(p))
	log.Info("Build pipeline")
	pipeline := &pipeline.Pipeline{
		Play:   p,
		Scheme: wf.Scheme,
	}
	pipeline.Parse(log)

	if err := wf.Add(pipeline.Create); err != nil {
		return err
	}
	return nil
}
