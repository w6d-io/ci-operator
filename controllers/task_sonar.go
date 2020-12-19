/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 29/11/2020
*/

package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

// SonarTask task struct for CI
type SonarTask struct {
	OwnerReferences []metav1.OwnerReference
	NamespacedName  types.NamespacedName
	Labels          map[string]string
	Steps           []tkn.Step
}

// Sonar create the sonar Tekton Task resource
func (wf *WFType) Sonar(ctx context.Context, p ci.Play, r *PlayReconciler) error {
	log := r.Log.WithName("Sonar").WithValues("task", ci.Sonar)
	// get the task
	log.V(1).Info("get task")
	// TODO get task from index in WFType
	steps, err := GetSteps(ctx, ci.Sonar, p, r)
	if err != nil {
		log.Error(err, "get steps failed")
		return err
	}
	if len(steps) == 0 {
		return fmt.Errorf("no step found for %s", ci.IntegrationTests)
	}
	sonar := &SonarTask{
		OwnerReferences: []metav1.OwnerReference{CIOwnerReference(p)},
		NamespacedName:  CxCINamespacedName(ci.Sonar.String(), p),
		Labels:          CxCILabels(p),
		Steps:           steps,
	}
	log.V(1).Info("add create in workflow")
	if err := wf.Add(sonar.Create); err != nil {
		return err
	}
	return nil
}

func (s *SonarTask) Create(ctx context.Context, r client.Client, log logr.Logger) error {
	log = log.WithName("Create").WithValues("task", ci.Sonar)
	log.V(1).Info("create")
	defaultMode := ci.FileMode0444
	// build Tekton Task resource
	taskResource := &tkn.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:            s.NamespacedName.Name,
			Namespace:       s.NamespacedName.Namespace,
			Annotations:     make(map[string]string),
			Labels:          s.Labels,
			OwnerReferences: s.OwnerReferences,
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
			Volumes: []corev1.Volume{
				{
					Name: "sqtoken",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName:  s.NamespacedName.Name,
							DefaultMode: &defaultMode,
						},
					},
				},
			},
		},
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
