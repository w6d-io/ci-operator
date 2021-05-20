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

// GitLeaksTask task struct for CI
type GitLeaksTask struct {
	Meta
}

// GitLeaks create git leaks tekton Task resource
func (t *Task) GitLeaks(ctx context.Context, logger logr.Logger) error {
	log := logger.WithName("GitLeaks").WithValues("task", ci.GitLeaks)
	// get the steps
	log.V(1).Info("build task")
	s := &Step{
		Index:    t.Index,
		PlaySpec: t.Play.Spec,
		Client:   t.Client,
		TaskType: ci.GitLeaks,
	}
	steps, err := s.GetSteps(ctx, logger)
	if err != nil {
		log.Error(err, "get steps failed")
		return err
	}
	if len(steps) == 0 {
		return fmt.Errorf("no step found for %s", s.TaskType)
	}
	task := t.Play.Spec.Tasks[t.Index][s.TaskType]
	if len(task.Variables) != 0 {
		for i := range steps {
			for key, val := range task.Variables {
				steps[i].Env = append(steps[i].Env, corev1.EnvVar{
					Name:  key,
					Value: val,
				})
			}
		}
	}
	gitLeaks := &GitLeaksTask{
		Meta{
			Steps:  steps,
			Play:   t.Play,
			Scheme: t.Scheme,
		},
	}

	log.V(1).Info("add create in workflow")
	return t.Add(gitLeaks.Create)
}

func (g *GitLeaksTask) Create(ctx context.Context, r client.Client, logger logr.Logger) error {
	log := logger.WithName("Create").WithValues("task", ci.GitLeaks)
	log.V(1).Info("creating")
	namespacedName := util.GetCINamespacedName(ci.GitLeaks.String(), g.Play)
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
