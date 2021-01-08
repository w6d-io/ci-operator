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

package task

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/util"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// UnitTestTask task struct for CI
type UnitTestTask struct {
	Meta
}

// UnitTest create the unitTest Tekton Task resource
func (t *Task) UnitTest(ctx context.Context, logger logr.Logger) error {
	log := logger.WithName("UnitTest").WithValues("task", ci.UnitTests)
	// get the task
	log.V(1).Info("get task")
	// TODO get task from index in WFType
	s := &Step{
		Index:    t.Index,
		PlaySpec: t.Play.Spec,
		Client:   t.Client,
		TaskType: ci.UnitTests,
	}
	steps, err := s.GetSteps(ctx, logger)
	if err != nil {
		log.Error(err, "get steps failed")
		return err
	}
	if len(steps) == 0 {
		return fmt.Errorf("no step found for %s", ci.UnitTests)
	}
	unittest := &UnitTestTask{
		Meta: Meta{
			Steps:  steps,
			Play:   t.Play,
			Scheme: t.Scheme,
		},
	}

	log.V(1).Info("add create in workflow")
	if err := t.Add(unittest.Create); err != nil {
		log.Error(err, "add function failed")
		return err
	}
	return nil
}

func (u *UnitTestTask) Create(ctx context.Context, r client.Client, log logr.Logger) error {
	log = log.WithName("Create").WithValues("task", ci.UnitTests)
	log.V(1).Info("creating")
	namespacedName := util.GetCINamespacedName(ci.UnitTests.String(), u.Play)

	// build Tekton Task resource
	taskResource := &tkn.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:        namespacedName.Name,
			Namespace:   namespacedName.Namespace,
			Annotations: make(map[string]string),
			Labels:      util.GetCILabels(u.Play),
		},
		Spec: tkn.TaskSpec{
			Workspaces: config.Workspaces(),
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
					Name: "script",
					Type: tkn.ParamTypeString,
				},
			},
			Steps: u.Steps,
		},
	}

	// set the current time in the resource annotations
	taskResource.Annotations[config.ScheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
	if err := controllerutil.SetControllerReference(u.Play, taskResource, u.Scheme); err != nil {
		return err
	}
	log.V(2).Info(fmt.Sprintf("task contains\n%v", util.GetObjectContain(taskResource)))
	if err := r.Create(ctx, taskResource); err != nil {
		return err
	}
	// All went well
	return nil
}
