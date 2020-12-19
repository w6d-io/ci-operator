/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 23/11/2020
*/

package controllers

import (
	"context"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"strings"
)

// GetStage build the tasks
func IsBuildStage(play ci.Play) bool {
	if strings.ToLower(play.Spec.Stack.Language) == "android" ||
		strings.ToLower(play.Spec.Stack.Language) == "ios" {
		return false
	}
	for _, t := range play.Spec.Tasks {
		for taskType := range t {
			if taskType == ci.Build {
				return true
			}
		}
	}
	return false
}

// SetTask executes the Task according the TaskType
//   - Build    => create a build Tekton Task
//   - UnitTest => create a unit test Tekton Task
//   - IntTest  => create a integration test Tekton Task
//   - Deploy   => create a deployment Tekton Task
//   - Clean    => create a cleaning Tekton Task
func (wf *WFType) SetTask(ctx context.Context, p ci.Play, r *PlayReconciler) error {
	log := r.Log.WithName("SetTask").WithValues("cx-namespace", InNamespace(p))
	for _, m := range p.Spec.Tasks {
		for name := range m {
			switch name {
			case ci.Build:
				// TODO call method build
				log.WithValues("task", ci.Build).V(1).Info("launch")
				if err := wf.Build(ctx, p, r); err != nil {
					return err
				}
			case ci.Sonar:
				// TODO call methods for unitTest
				log.WithValues("task", ci.Sonar).V(1).Info("launch")
				if err := wf.Sonar(ctx, p, r); err != nil {
					return err
				}
			case ci.UnitTests:
				// TODO call methods for unitTest
				log.WithValues("task", ci.UnitTests).V(1).Info("launch")
				if err := wf.UnitTest(ctx, p, r); err != nil {
					return err
				}
			case ci.IntegrationTests:
				// TODO call methods for integrationTest
				log.WithValues("task", ci.IntegrationTests).V(1).Info("launch")
				if err := wf.IntTest(ctx, p, r); err != nil {
					return err
				}
			case ci.Deploy:
				// TODO call methods for deploy
				log.WithValues("task", ci.Deploy).V(1).Info("launch")
				if err := wf.Deploy(ctx, p, r); err != nil {
					return err
				}
			case ci.Clean:
				// TODO call methods for clean
				log.WithValues("task", ci.Clean).V(1).Info("launch")
				if err := wf.Clean(ctx, p, r); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
