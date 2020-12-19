/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 22/11/2020
*/

package controllers

import (
	"context"
	"github.com/go-logr/logr"
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/tektoncd/pipeline/pkg/apis/resource/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
)

// CreateCI takes a Play struct to create tekton Pipeline
func CreateCI(ctx context.Context, p ci.Play, r *PlayReconciler) error {

	log := r.Log.WithName("CreateCI").WithValues("cx-namespace", InNamespace(p))
	p.Status.State = ci.Creating
	if err := r.Status().Update(ctx, &p); err != nil {
		return err
	}

	wf := New()
	if err := wf.SetGitCreate(p, r); err != nil {
		return err
	}
	if err := wf.SetImageCreate(p, r); err != nil {
		return err
	}
	if err := wf.SetTask(ctx, p, r); err != nil {
		return err
	}
	if err := wf.SetPipeline(p, r); err != nil {
		return err
	}
	if err := wf.SetPipelineRun(p, r); err != nil {
		return err
	}
	log.Info("Launch creation")
	if err := wf.Run(ctx, r, log); err != nil {
		log.Error(err, "CI creation failed")
		// TODO add rollback ( delete resource created before )
		return err
	}
	return nil
}

// TODO Add is ci exists function / method

// Run executes Create methods in WFType
func (wf *WFType) Run(ctx context.Context, r client.Client, log logr.Logger) error {
	for _, c := range wf.Creates {
		if err := c(ctx, r, log); err != nil {
			return err
		}
	}
	return nil
}

// CleanCI remove tekton resource
func CleanCI(ctx context.Context, p ci.Play, r *PlayReconciler) {
	log := r.Log.WithName("CleanCI").WithValues("cx-namespace", InNamespace(p))

	if r == nil {
		return
	}
	log.V(1).Info("get pipeline resource")
	//var prs v1alpha1.PipelineResourceList
	var labels client.MatchingLabels
	labels = CxCILabels(p)
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(v1alpha1.SchemeGroupVersion.WithKind("PipelineResource"))
	log.V(2).Info("delete pipeline resource")
	if err := r.DeleteAllOf(ctx, u, InNamespace(p), labels); err != nil {
		log.Error(err, "delete pipeline resource failed")
	}

	u.SetGroupVersionKind(tkn.SchemeGroupVersion.WithKind("Task"))
	log.V(2).Info("delete task")
	if err := r.DeleteAllOf(ctx, u, InNamespace(p), labels); err != nil {
		log.Error(err, "delete tasks failed")
	}

	u.SetGroupVersionKind(tkn.SchemeGroupVersion.WithKind("Pipeline"))
	log.V(2).Info("delete pipeline")
	if err := r.DeleteAllOf(ctx, u, InNamespace(p), labels); err != nil {
		log.Error(err, "delete pipeline failed")
	}

	u.SetGroupVersionKind(tkn.SchemeGroupVersion.WithKind("PipelineRun"))
	log.V(2).Info("delete pipeline run")
	if err := r.DeleteAllOf(ctx, u, InNamespace(p), labels); err != nil {
		log.Error(err, "delete pipeline run failed")
	}
}
