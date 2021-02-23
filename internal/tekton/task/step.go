/*
Copyright 2020 WILDCARD

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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

// TODO those are duplicates entries for breaking cycle import find a way to remove it
const (
	// filename use with s3cmd
	MinIOSecretKey = ".s3cfg"
	// Prefix use for name of resource
	MinIOPrefixSecret = "minio"
)

// Step structure for GetStep and FilteredStep
type Step struct {
	Index    int
	PlaySpec ci.PlaySpec
	Client   client.Client
	TaskType ci.TaskType
}

// GetSteps return the list of step according the task
func (s *Step) GetSteps(ctx context.Context, logger logr.Logger) ([]tkn.Step, error) {
	log := logger.WithName("GetSteps").WithValues("task", s.TaskType)
	// get Step by annotation
	var steplist ci.StepList
	//var opts []client.ListOption

	err := s.Client.List(ctx, &steplist)
	if err != nil {
		return nil, err
	}
	log.WithValues("nbr", len(steplist.Items)).V(2).Info("List return")
	sortedSteps := s.FilteredSteps(logger, steplist.Items, s.TaskType == ci.UnitTests || s.TaskType == ci.IntegrationTests)
	log.WithValues("nbr", len(sortedSteps)).V(2).Info("Filtered list return")
	if len(sortedSteps) == 0 {
		log.Error(errors.New("get steps error"), "list empty")
		return []tkn.Step{}, nil
	}
	sort.Sort(&sortedSteps)
	var steps []tkn.Step
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
func (s *Step) FilteredSteps(logger logr.Logger, steps ci.Steps, isTest bool) ci.Steps {
	filteredSteps := ci.Steps{}
	log := logger.WithName("FilteredSteps").WithValues("task", s.TaskType, "stack", s.PlaySpec.Stack,
		"ops-namespace", config.GetNamespace())
	log.V(2).Info("filtering")
	_, mongoOK := s.PlaySpec.Dependencies[ci.MongoDB]
	_, postgresOK := s.PlaySpec.Dependencies[ci.Postgresql]
	_, mariaDBOK := s.PlaySpec.Dependencies[ci.MariaDB]
	task := s.PlaySpec.Tasks[s.Index][s.TaskType]

	for _, step := range steps {
		if config.GetNamespace() != "" && step.Namespace != config.GetNamespace() {
			continue
		}
		if (mongoOK || postgresOK || mariaDBOK) && step.Annotations[ci.AnnotationTask] == s.TaskType.String() &&
			(step.Annotations[ci.AnnotationLanguage] == ci.MongoDB.String() ||
				step.Annotations[ci.AnnotationLanguage] == ci.Postgresql.String()) {
			filteredSteps = append(filteredSteps, step)
			continue
		}
		log.WithValues("package", step.Annotations[ci.AnnotationPackage],
			"task", step.Annotations[ci.AnnotationTask],
			"language", step.Annotations[ci.AnnotationLanguage]).V(2).Info("annotations")
		if isTest {
			if (len(task.Script) == 0) && (step.Annotations[ci.AnnotationPackage] != s.PlaySpec.Stack.Package) {
				continue
			}
			if (len(task.Script) != 0) && (step.Annotations[ci.AnnotationPackage] != "custom") {
				continue
			}
		}
		if step.Annotations[ci.AnnotationTask] != s.TaskType.String() {
			continue
		}
		if step.Annotations[ci.AnnotationLanguage] != s.PlaySpec.Stack.Language {
			continue
		}
		filteredSteps = append(filteredSteps, step)
	}
	if len(filteredSteps) == 0 {
		filteredSteps = s.GetGenericSteps(logger, steps)
	}
	return filteredSteps
}

// GetGenericStep returns the steps bind with the task type
func (s *Step) GetGenericSteps(logger logr.Logger, steps ci.Steps) ci.Steps {
	log := logger.WithName("GetGenericSteps")
	data := ci.Steps{}
	log.V(1).Info("get generic steps")

	for _, step := range steps {
		if config.GetNamespace() != "" && step.Namespace != config.GetNamespace() {
			continue
		}
		if step.Annotations[ci.AnnotationKind] != "generic" {
			continue
		}
		if step.Annotations[ci.AnnotationTask] != s.TaskType.String() {
			continue
		}
		data = append(data, step)
	}
	return data
}
