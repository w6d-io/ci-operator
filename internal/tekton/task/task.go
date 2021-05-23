/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 17/05/2021
*/

package task

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// GenericTask task struct for CI
type GenericTask struct {
	Meta
	TaskType ci.TaskType
	Params   []ci.ParamSpec
}

// Generic create git leaks tekton Task resource
func (t *Task) Generic(ctx context.Context, taskType ci.TaskType, logger logr.Logger) error {
	log := logger.WithName("Generic").WithValues("task", taskType)
	// get the steps
	log.V(1).Info("build task")
	s := &Step{
		Index:    t.Index,
		PlaySpec: t.Play.Spec,
		Client:   t.Client,
		TaskType: taskType,
	}
	steps, params, err := s.GetSteps(ctx, logger)
	if err != nil {
		log.Error(err, "get steps failed")
		return err
	}
	t.Params[string(taskType)] = params
	if len(steps) == 0 {
		return fmt.Errorf("no step found for %s", taskType)
	}
	task := t.Play.Spec.Tasks[t.Index][taskType]
	for i := range steps {
		if len(task.Variables) != 0 {
			for key, val := range task.Variables {
				steps[i].Env = append(steps[i].Env, corev1.EnvVar{
					Name:  key,
					Value: val,
				})
			}
		}
		steps[i].Env = append(steps[i].Env, BuildAndGetPredefinedEnv(t.Play)...)
	}
	g := &GenericTask{
		Meta{
			Steps:  steps,
			Play:   t.Play,
			Scheme: t.Scheme,
		},
		taskType,
		params,
	}

	log.V(1).Info("add create in workflow")
	return t.Add(g.Create)
}

func (g *GenericTask) Create(ctx context.Context, r client.Client, logger logr.Logger) error {
	log := logger.WithName("Create").WithValues("task", g.TaskType)
	log.V(1).Info("creating")
	namespacedName := util.GetCINamespacedName(g.TaskType.String(), g.Play)
	// build Tekton Task resource
	resource := &tkn.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:        namespacedName.Name,
			Namespace:   namespacedName.Namespace,
			Annotations: make(map[string]string),
			Labels:      util.GetCILabels(g.Play),
		},
		Spec: tkn.TaskSpec{
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
			Steps:      g.Steps,
			Workspaces: config.Workspaces(),
		},
	}
	for _, param := range g.Params {
		resource.Spec.Params = append(resource.Spec.Params, tkn.ParamSpec{
			Name: param.Name,
			Type: param.Type,
		})
	}
	// set the current time in the annotations
	resource.Annotations[config.ScheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
	if err := controllerutil.SetControllerReference(g.Play, resource, g.Scheme); err != nil {
		return err
	}
	log.V(1).Info(resource.Kind, "content", fmt.Sprintf("%v",
		util.GetObjectContain(resource)))
	if err := r.Create(ctx, resource); err != nil {
		return err
	}
	return nil
}
