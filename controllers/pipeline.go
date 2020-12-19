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

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (wf *WFType) SetPipeline(p ci.Play, r *PlayReconciler) error {
	log := r.Log.WithName("SetPipeline").WithValues("cx-namespace", InNamespace(p))
	log.Info("Build pipeline")
	pipeline := &Pipeline{
		OwnerReferences: []metav1.OwnerReference{CIOwnerReference(p)},
		NamespacedName:  CxCINamespacedName("pipeline", p),
		Labels:          CxCILabels(p),
		Resources: []tkn.PipelineDeclaredResource{
			{
				Name: ci.ResourceGit,
				Type: tkn.PipelineResourceTypeGit,
			},
		},
	}
	for _, wks := range Cfg.Workspaces {
		pipeline.Workspaces = append(pipeline.Workspaces, tkn.PipelineWorkspaceDeclaration{
			Name: wks.Name,
		})
	}
	if IsBuildStage(p) {
		pipeline.Resources = append(pipeline.Resources, tkn.PipelineDeclaredResource{
			Name: ci.ResourceImage,
			Type: tkn.PipelineResourceTypeImage,
		})
	}
	for _, m := range p.Spec.Tasks {
		for name := range m {
			switch name {
			case ci.Build:
				if err := pipeline.SetPipelineBuild(p, r); err != nil {
					return err
				}
			case ci.Sonar:
				if err := pipeline.SetPipelineSonar(p, r); err != nil {
					return err
				}
			case ci.UnitTests:
				if err := pipeline.SetPipelineUnitTest(p, r); err != nil {
					return err
				}
			case ci.IntegrationTests:
				if err := pipeline.SetPipelineIntTest(p, r); err != nil {
					return err
				}
			case ci.Deploy:
				if err := pipeline.SetPipelineDeploy(p, r); err != nil {
					return err
				}
			case ci.Clean:
				if err := pipeline.SetPipelineClean(p, r); err != nil {
					return err
				}
			}
		}
	}

	if err := wf.Add(pipeline.Create); err != nil {
		return err
	}
	return nil
}

// Create
func (p *Pipeline) Create(ctx context.Context, r client.Client, log logr.Logger) error {
	log = log.WithName("Create").WithValues("action", "pipeline")
	log.V(1).Info("creating")

	pipelineResource := &tkn.Pipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:            p.NamespacedName.Name,
			Namespace:       p.NamespacedName.Namespace,
			Annotations:     make(map[string]string),
			Labels:          p.Labels,
			OwnerReferences: p.OwnerReferences,
		},
		Spec: tkn.PipelineSpec{
			Params:     p.Params,
			Resources:  p.Resources,
			Tasks:      p.Tasks,
			Workspaces: p.Workspaces,
		},
	}

	log.V(1).Info(fmt.Sprintf("pipeline contains\n%v", GetObjectContain(pipelineResource)))
	pipelineResource.Annotations[scheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
	if err := r.Create(ctx, pipelineResource); err != nil {
		return err
	}
	// All went well
	return nil
}

func getWorkspacePipelineTaskBinding() []tkn.WorkspacePipelineTaskBinding {
	var ws []tkn.WorkspacePipelineTaskBinding

	for _, wks := range Cfg.Workspaces {
		wb := tkn.WorkspacePipelineTaskBinding{
			Name:      wks.Name,
			Workspace: wks.Name,
			SubPath:   "/workspaces/" + wks.Name,
		}
		if wks.MountPath != "" {
			wb.SubPath = wks.MountPath
		}
		ws = append(ws, wb)
	}

	return ws
}

// GetWorkspacePath return path from workspace
func GetWorkspacePath(name string, wks []tkn.WorkspaceDeclaration) string {
	for _, wk := range wks {
		if wk.Name == name {
			subPath := "/workspaces/" + wk.Name
			if wk.MountPath != "" {
				subPath = wk.MountPath
			}
			return subPath
		}
	}
	return ""
}

// return a pointer on task resource for pipeline
func getPipelineTaskResources(build bool) *tkn.PipelineTaskResources {
	ptr := &tkn.PipelineTaskResources{}

	ptr.Inputs = []tkn.PipelineTaskInputResource{
		{
			Name:     ci.ResourceGit,
			Resource: ci.ResourceGit,
		},
	}
	if build {
		ptr.Outputs = []tkn.PipelineTaskOutputResource{
			{
				Name:     ci.ResourceImage,
				Resource: ci.ResourceImage,
			},
		}
	}
	return ptr
}

func getParamString(name string) tkn.Param {
	return tkn.Param{
		Name: name,
		Value: tkn.ArrayOrString{
			Type:      tkn.ParamTypeString,
			StringVal: "$(params." + name + ")",
		},
	}
}
