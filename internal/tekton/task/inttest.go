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
Created on 24/11/2020
*/

package task

import (
	"context"
	"fmt"
	"time"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// IntTestTask task struct for CI
type IntTestTask struct {
	Meta
}

// IntTest create the intTest Tekton Task resource
func (t *Task) IntTest(ctx context.Context, logger logr.Logger) error {
	log := logger.WithName("IntTest").WithValues("task", ci.IntegrationTests)
	// get the task
	log.V(1).Info("get task")
	s := &Step{
		Index:    t.Index,
		PlaySpec: t.Play.Spec,
		Client:   t.Client,
		TaskType: ci.IntegrationTests,
	}
	steps, _, _, err := s.GetSteps(ctx, logger)
	if err != nil {
		log.Error(err, "get steps failed")
		return err
	}
	if len(steps) == 0 {
		return fmt.Errorf("no step found for %s", s.TaskType)
	}
	task := t.Play.Spec.Tasks[t.Index][s.TaskType]
	for i := range steps {
		for key, val := range task.Variables {
			steps[i].Env = append(steps[i].Env, corev1.EnvVar{
				Name:  key,
				Value: val,
			})
		}
		steps[i].Env = append(steps[i].Env, BuildAndGetPredefinedEnv(t.Play)...)
	}
	inttest := &IntTestTask{
		Meta: Meta{
			Steps:  steps,
			Play:   t.Play,
			Scheme: t.Scheme,
		},
	}

	log.V(1).Info("add create in workflow")
	return t.Add(inttest.Create)
}

func (u *IntTestTask) Create(ctx context.Context, r client.Client, log logr.Logger) error {
	log = log.WithName("Create").WithValues("task", ci.IntegrationTests)
	log.V(1).Info("creating")
	namespacedName := util.GetCINamespacedName(ci.IntegrationTests.String(), u.Play)
	// build Tekton Task resource
	resource := &tkn.Task{
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
					Name: "IMAGE",
					Type: tkn.ParamTypeString,
				},
				{
					Name: "script",
					Type: tkn.ParamTypeString,
				},
			},
			Steps: u.Steps,
		},
	}

	// set the current time in the resource annotations
	resource.Annotations[config.ScheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
	if err := controllerutil.SetControllerReference(u.Play, resource, u.Scheme); err != nil {
		return err
	}
	log.V(1).Info(resource.Kind, "content", fmt.Sprintf("%v",
		util.GetObjectContain(resource)))
	if err := r.Create(ctx, resource); err != nil {
		return err
	}
	// All went well
	return nil
}
