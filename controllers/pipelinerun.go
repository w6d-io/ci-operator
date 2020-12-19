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

package controllers

import (
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/controllers/pipelinerun"
)

func (wf *WFType) SetPipelineRun(play ci.Play, r *PlayReconciler) error {
	log := r.Log.WithName("SetPipelineRun").WithValues("cx-namespace", InNamespace(play))
	log.Info("Build pipeline run")

	pipelineRun := &pipelinerun.PipelineRun{
		GetNamespacedName:     CxCINamespacedName,
		GetLabels:             CxCILabels,
		GetOwnerReferences:    CIOwnerReference,
		GetObjectContain:      GetObjectContain,
		GetDockerImage:        CxDockerImageName,
		DeployNamespacedName:  CxDeployNamespacedName,
		IsBuildStage:          IsBuildStage,
		GetWorkspace:          GetWorkspacePath,
		ServiceAccount:        Cfg.ServiceAccountName,
		WorkspacesDeclaration: Cfg.Workspaces,
		Volume: pipelinerun.Volume{
			Size: Cfg.Volume.Size,
			Mode: Cfg.Volume.Mode,
		},
	}

	if err := wf.Add(pipelineRun.Parse(play, r.Log)); err != nil {
		return err
	}
	return nil
}
