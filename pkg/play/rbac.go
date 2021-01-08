/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 07/01/2021
*/

package play

import (
	"github.com/go-logr/logr"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal"
	"github.com/w6d-io/ci-operator/internal/k8s/rbac"
	"github.com/w6d-io/ci-operator/internal/util"
)

func (wf *WFType) Rbac(p *ci.Play, logger logr.Logger) error {
	log := logger.WithName("GitSecret").WithValues("cx-namespace", util.InNamespace(p))
	log.V(1).Info("Build git sa")

	resource := &rbac.Rbac{
		WorkFlowStruct: internal.WorkFlowStruct{
			Play:   p,
			Scheme: wf.Scheme,
		},
	}

	if err := wf.Add(resource.Create); err != nil {
		return err
	}
	return nil
}
