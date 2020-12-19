/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 25/11/2020
*/

package controllers

import (
	"context"
	"errors"
	"github.com/go-logr/logr"
	"sort"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
)

// GetSteps return the list of step according the task
func GetSteps(ctx context.Context, taskType ci.TaskType, p ci.Play, r *PlayReconciler) ([]tkn.Step, error) {
	log := r.Log.WithName("GetSteps").WithValues("task", taskType)
	// get Step by annotation
	var steplist ci.StepList
	//var opts []client.ListOption

	//opts = append(opts, client.MatchingFields{"metadata.annotations." + ci.AnnotationTask: string(taskType)})
	//opts = append(opts, client.MatchingFields{"metadata.annotations." + ci.AnnotationLanguage: scope.Language})
	//
	//if taskType == ci.UnitTests || taskType == ci.IntegrationTests {
	//	opts = append(opts, client.MatchingFields{"metadata.annotations." + ci.AnnotationPackage: scope.Package})
	//}
	err := r.List(ctx, &steplist)
	if err != nil {
		return nil, err
	}
	log.WithValues("nbr", len(steplist.Items)).V(2).Info("List return")
	sortedSteps := FilteredSteps(r.Log, steplist.Items, p.Spec.Stack, taskType, taskType == ci.UnitTests || taskType == ci.IntegrationTests)
	log.WithValues("nbr", len(sortedSteps)).V(2).Info("Filtered list return")
	if len(sortedSteps) == 0 {
		log.Error(errors.New("get steps error"), "list empty")
		return []tkn.Step{}, nil
	}
	sort.Sort(&sortedSteps)
	var steps []tkn.Step
	// TODO get step by annotation in Step kind
	for _, step := range sortedSteps {
		steps = append(steps, tkn.Step{
			Container: step.Step.Container,
			Script:    step.Step.Script,
			Timeout:   step.Step.Timeout,
		})
	}
	return steps, nil
}

// FilteredSteps return a ci.Steps filtered by annotation
func FilteredSteps(log logr.Logger, steps ci.Steps, stack ci.Stack, taskType ci.TaskType, isTest bool) ci.Steps {
	filteredSteps := ci.Steps{}
	log = log.WithName("FilteredSteps").WithValues("task", taskType, "stack", stack)
	for _, step := range steps {
		log.WithValues("package", step.Annotations[ci.AnnotationPackage],
			"task", step.Annotations[ci.AnnotationTask],
			"language", step.Annotations[ci.AnnotationLanguage]).V(2).Info("annotations")
		if isTest && (step.Annotations[ci.AnnotationPackage] != stack.Package && step.Annotations[ci.AnnotationPackage] != "custom") {
			continue
		}
		if step.Annotations[ci.AnnotationTask] != string(taskType) {
			continue
		}
		if step.Annotations[ci.AnnotationLanguage] != stack.Language {
			continue
		}
		filteredSteps = append(filteredSteps, step)
	}
	return filteredSteps
}
