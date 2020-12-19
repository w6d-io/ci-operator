/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 24/11/2020
*/

package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

// BuildTask task struct fo CI
type BuildTask struct {
	OwnerReferences []metav1.OwnerReference
	NamespacedName  types.NamespacedName
	Labels          map[string]string
	Steps           []tkn.Step
	BuildDocker     bool
}

// Build create the build Tekton Task resource
func (wf *WFType) Build(ctx context.Context, p ci.Play, r *PlayReconciler) error {
	log := r.Log.WithName("Build").WithValues("task", ci.Build)
	log.V(1).Info("get steps")
	steps, err := GetSteps(ctx, ci.Build, p, r)
	if err != nil {
		log.Error(err, "get steps failed")
		return err
	}
	if len(steps) == 0 {
		return fmt.Errorf("no step found for %s", ci.Build)
	}
	build := &BuildTask{
		OwnerReferences: []metav1.OwnerReference{CIOwnerReference(p)},
		NamespacedName:  CxCINamespacedName(string(ci.Build), p),
		Labels:          CxCILabels(p),
		Steps:           steps,
		BuildDocker:     IsBuildStage(p),
	}
	log.V(1).Info("add create in workflow")
	if err := wf.Add(build.Create); err != nil {
		log.Error(err, "add function failed")
		return err
	}
	return nil
}

func (b *BuildTask) Create(ctx context.Context, r client.Client, log logr.Logger) error {
	// build Tekton Task resource
	log = log.WithName("Create").WithValues("task", ci.Build)
	log.V(1).Info("creating")

	taskResource := &tkn.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:            b.NamespacedName.Name,
			Namespace:       b.NamespacedName.Namespace,
			Annotations:     make(map[string]string),
			Labels:          b.Labels,
			OwnerReferences: b.OwnerReferences,
		},
		Spec: tkn.TaskSpec{
			Workspaces: Cfg.Workspaces,
			Resources: &tkn.TaskResources{
				Inputs: []tkn.TaskResource{
					{
						ResourceDeclaration: tkn.ResourceDeclaration{
							Name: ci.ResourceGit,
							Type: tkn.PipelineResourceTypeGit,
						},
					},
				},
			},
			Params: []tkn.ParamSpec{
				{
					Name: "flags",
					Type: tkn.ParamTypeArray,
				},
			},
			Steps: b.Steps,
		},
	}
	if b.BuildDocker {
		taskResource.Spec.Resources.Outputs = []tkn.TaskResource{
			{
				ResourceDeclaration: tkn.ResourceDeclaration{
					Name: ci.ResourceImage,
					Type: tkn.PipelineResourceTypeImage,
				},
			},
		}
		taskResource.Spec.Params = append(taskResource.Spec.Params, []tkn.ParamSpec{
			{
				Name: "DOCKERFILE",
				Type: tkn.ParamTypeString,
			},
			{
				Name: "IMAGE",
				Type: tkn.ParamTypeString,
			},
			{
				Name: "CONTEXT",
				Type: tkn.ParamTypeString,
			},
		}...)
	}

	log.V(2).Info(fmt.Sprintf("task contains\n%v", GetObjectContain(taskResource)))
	// set the current time in the resource annotations
	taskResource.Annotations[scheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
	if err := r.Create(ctx, taskResource); err != nil {
		return err
	}
	// All went well
	return nil
}
