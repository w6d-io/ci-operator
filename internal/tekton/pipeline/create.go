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

package pipeline

import (
	"context"
	"fmt"
	"time"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const Prefix string = "pipeline"

func (p *Pipeline) Parse(logger logr.Logger) error {
	p.Workspaces = append(p.Workspaces, tkn.PipelineWorkspaceDeclaration{
		Name: config.Volume().Name,
	})
	for pos, m := range p.Play.Spec.Tasks {
		for name := range m {
			switch name {
			case ci.Build:
				p.Pos = pos
				if err := p.SetPipelineBuild(p.Play, logger); err != nil {
					return err
				}
			case ci.Sonar:
				if err := p.SetPipelineSonar(p.Play, logger); err != nil {
					return err
				}
			case ci.UnitTests:
				if err := p.SetPipelineUnitTest(p.Play, logger); err != nil {
					return err
				}
			case ci.IntegrationTests:
				if err := p.SetPipelineIntTest(p.Play, logger); err != nil {
					return err
				}
			case ci.Deploy:
				if err := p.SetPipelineDeploy(p.Play, logger); err != nil {
					return err
				}
			case ci.Clean:
				if err := p.SetPipelineClean(p.Play, logger); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Create
func (p *Pipeline) Create(ctx context.Context, r client.Client, log logr.Logger) error {
	log = log.WithName("Create").WithValues("action", Prefix)
	log.V(1).Info("creating")

	namespacedName := util.GetCINamespacedName(Prefix, p.Play)
	resources := []tkn.PipelineDeclaredResource{
		{
			Name: ci.ResourceGit,
			Type: tkn.PipelineResourceTypeGit,
		},
	}
	if util.IsBuildStage(p.Play) {
		resources = append(resources, tkn.PipelineDeclaredResource{
			Name: ci.ResourceImage,
			Type: tkn.PipelineResourceTypeImage,
		})
	}
	pipelineResource := &tkn.Pipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:        namespacedName.Name,
			Namespace:   namespacedName.Namespace,
			Annotations: make(map[string]string),
			Labels:      util.GetCILabels(p.Play),
		},
		Spec: tkn.PipelineSpec{
			Params:     p.Params,
			Resources:  resources,
			Tasks:      p.Tasks,
			Workspaces: p.Workspaces,
		},
	}

	pipelineResource.Annotations[config.ScheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
	if err := controllerutil.SetControllerReference(p.Play, pipelineResource, p.Scheme); err != nil {
		return err
	}
	log.V(1).Info(fmt.Sprintf("pipeline contains\n%v", util.GetObjectContain(pipelineResource)))
	if err := r.Create(ctx, pipelineResource); err != nil {
		return err
	}
	// All went well
	return nil
}

func getWorkspacePipelineTaskBinding() []tkn.WorkspacePipelineTaskBinding {
	var ws []tkn.WorkspacePipelineTaskBinding

	for _, wks := range config.Workspaces() {
		wb := tkn.WorkspacePipelineTaskBinding{
			Name:      wks.Name,
			Workspace: config.Volume().Name,
			SubPath:   wks.Name,
		}
		ws = append(ws, wb)
	}

	return ws
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
