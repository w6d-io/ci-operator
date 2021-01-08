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

package task

import (
	"context"
	"errors"
	"sort"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"

	"github.com/go-logr/logr"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline"
	"github.com/w6d-io/ci-operator/internal/config"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TODO they are duplicates entries for breaking cycle import
const (
	// filename use with s3cmd
	MinIOSecretKey = ".s3cfg"
	// Prefix use for name of resource
	MinIOPrefixSecret = "minio"
)

// GetSteps return the list of step according the task
func GetSteps(ctx context.Context, taskType ci.TaskType, p *ci.Play, logger logr.Logger, r client.Client) ([]tkn.Step, error) {
	log := logger.WithName("GetSteps").WithValues("task", taskType)
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
	sortedSteps := FilteredSteps(logger, steplist.Items, p.Spec, taskType, taskType == ci.UnitTests || taskType == ci.IntegrationTests)
	log.WithValues("nbr", len(sortedSteps)).V(2).Info("Filtered list return")
	if len(sortedSteps) == 0 {
		log.Error(errors.New("get steps error"), "list empty")
		return []tkn.Step{}, nil
	}
	sort.Sort(&sortedSteps)
	var steps []tkn.Step
	// TODO get step by annotation in Step kind
	for _, step := range sortedSteps {
		newStep := tkn.Step{
			Container: step.Step.Container,
			Script:    step.Step.Script,
			Timeout:   step.Step.Timeout,
		}
		if config.GetMinio().Host != "" {
			vol := corev1.VolumeMount{
				MountPath: pipeline.HomeDir + "/" + MinIOSecretKey,
				Name:      MinIOPrefixSecret,
				SubPath:   MinIOSecretKey,
			}
			newStep.Container.VolumeMounts = append(newStep.Container.VolumeMounts, vol)
		}
		steps = append(steps, newStep)
	}
	return steps, nil
}

// FilteredSteps return a ci.Steps filtered by annotation
func FilteredSteps(log logr.Logger, steps ci.Steps, spec ci.PlaySpec, taskType ci.TaskType, isTest bool) ci.Steps {
	filteredSteps := ci.Steps{}
	log = log.WithName("FilteredSteps").WithValues("task", taskType, "stack", spec.Stack)
	_, mongoOK := spec.Dependencies[ci.MongoDB]
	_, postgresOK := spec.Dependencies[ci.Postgresql]
	_, mariaDBOK := spec.Dependencies[ci.MariaDB]

	for _, step := range steps {
		if (mongoOK || postgresOK || mariaDBOK) && step.Annotations[ci.AnnotationTask] == taskType.String() &&
			(step.Annotations[ci.AnnotationLanguage] == ci.MongoDB.String() ||
				step.Annotations[ci.AnnotationLanguage] == ci.Postgresql.String()) {
			filteredSteps = append(filteredSteps, step)
			continue
		}
		log.WithValues("package", step.Annotations[ci.AnnotationPackage],
			"task", step.Annotations[ci.AnnotationTask],
			"language", step.Annotations[ci.AnnotationLanguage]).V(2).Info("annotations")
		if isTest && (step.Annotations[ci.AnnotationPackage] != spec.Stack.Package) {
			continue
		}
		if step.Annotations[ci.AnnotationTask] != string(taskType) {
			continue
		}
		if step.Annotations[ci.AnnotationLanguage] != spec.Stack.Language {
			continue
		}
		filteredSteps = append(filteredSteps, step)
	}
	return filteredSteps
}
