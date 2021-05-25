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
	"github.com/w6d-io/ci-operator/internal/k8s/secrets"
	"path"
	"sort"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"

	"github.com/go-logr/logr"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline"
	"github.com/w6d-io/ci-operator/internal/config"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Step structure for GetStep and FilteredStep
type Step struct {
	Index    int
	PlaySpec ci.PlaySpec
	Client   client.Client
	TaskType ci.TaskType
}

// GetSteps return the list of step according the task
func (s *Step) GetSteps(ctx context.Context, logger logr.Logger) (steps []tkn.Step, params []ci.ParamSpec, sidecars []tkn.Sidecar, err error) {
	logger = logger.WithValues("task", s.TaskType)
	log := logger.WithName("GetSteps")
	// get Step by annotation
	var steplist ci.StepList
	//var opts []client.ListOption

	err = s.Client.List(ctx, &steplist)
	if err != nil {
		return
	}
	log.WithValues("nbr", len(steplist.Items)).V(1).Info("List return")
	if len(s.PlaySpec.Tasks) < s.Index+1 {
		err = errors.New("no such task")
		return
	}
	sortedSteps := s.FilteredSteps(logger, steplist.Items, s.TaskType == ci.UnitTests ||
		s.TaskType == ci.IntegrationTests || s.TaskType == ci.E2ETests)
	log.WithValues("nbr", len(sortedSteps)).V(1).Info("Filtered list return")
	if len(sortedSteps) == 0 {
		log.Error(errors.New("get steps error"), "list empty")
		return
	}
	sort.Sort(&sortedSteps)
	for _, step := range sortedSteps {
		newStep := tkn.Step{
			Container: step.Step.Container,
			Script:    step.Step.Script,
			Timeout:   step.Step.Timeout,
		}
		params = append(params, step.Params...)
		sidecars = append(sidecars, step.Sidecar...)
		if config.GetMinio().Host != "" {
			vol := corev1.VolumeMount{
				MountPath: path.Join(pipeline.HomeDir, secrets.MinIOSecretKey),
				Name:      secrets.MinIOPrefixSecret,
				SubPath:   secrets.MinIOSecretKey,
			}
			newStep.Container.VolumeMounts = append(newStep.Container.VolumeMounts, vol)
		}
		var okVault, okSecret bool
		if s.PlaySpec.Vault != nil {
			_, okVault = s.PlaySpec.Vault.Secrets[secrets.KubeConfigKey]
		}
		_, okSecret = s.PlaySpec.Secret[secrets.KubeConfigKey]
		if okVault || okSecret {
			vol := corev1.VolumeMount{
				Name:      secrets.KubeConfigPrefix,
				MountPath: pipeline.HomeDir + "/.kube/config",
				SubPath:   "config",
			}
			newStep.Container.VolumeMounts = append(newStep.Container.VolumeMounts, vol)
		}
		steps = append(steps, newStep)
	}
	return
}

// GetParams return the params from Step
func (s *Step) GetParams(params []ci.ParamSpec) (r []tkn.ParamSpec) {
	for i := range params {
		r = append(r, tkn.ParamSpec{
			Name:        params[i].Name,
			Type:        params[i].Type,
			Description: params[i].Description,
			Default:     params[i].Default,
		})
	}
	return
}

// FilteredSteps return a ci.Steps filtered by annotation
func (s *Step) FilteredSteps(logger logr.Logger, steps ci.Steps, isTest bool) ci.Steps {
	filteredSteps := ci.Steps{}
	log := logger.WithName("FilteredSteps").WithValues("stack", s.PlaySpec.Stack,
		"ops-namespace", config.GetNamespace())
	log.V(1).Info("filtering")
	task := s.PlaySpec.Tasks[s.Index][s.TaskType]

	for _, step := range steps {
		if config.GetNamespace() != "" && step.Namespace != config.GetNamespace() {
			continue
		}
		log.WithValues("step_package", step.Annotations[ci.AnnotationPackage],
			"step_task", step.Annotations[ci.AnnotationTask],
			"step_language", step.Annotations[ci.AnnotationLanguage]).V(1).Info("annotations")
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

// GetGenericSteps returns the steps bind with the task type
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
