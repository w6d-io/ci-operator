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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// BuildTask task struct fo CI
type BuildTask struct {
	Meta
	BuildDocker bool
}

// Build create the build Tekton Task resource
func (t *Task) Build(ctx context.Context, logger logr.Logger) error {
	log := logger.WithName("Build").WithValues("task", ci.Build)
	log.V(1).Info("get steps")
	s := &Step{
		Index:    t.Index,
		PlaySpec: t.Play.Spec,
		Client:   t.Client,
		TaskType: ci.Build,
	}
	steps, err := s.GetSteps(ctx, logger)
	if err != nil {
		log.Error(err, "get steps failed")
		return err
	}
	if len(steps) == 0 {
		return fmt.Errorf("no step found for %s", ci.Build)
	}
	build := &BuildTask{
		Meta: Meta{
			Steps:  steps,
			Play:   t.Play,
			Scheme: t.Scheme,
		},
		BuildDocker: util.IsBuildStage(t.Play),
	}
	log.V(1).Info("add create in workflow")
	if err := t.Add(build.Create); err != nil {
		log.Error(err, "add function failed")
		return err
	}
	return nil
}

func (b *BuildTask) Create(ctx context.Context, r client.Client, log logr.Logger) error {
	// build Tekton Task resource
	log = log.WithName("Create").WithValues("task", ci.Build)
	log.V(1).Info("creating")
	namespacedName := util.GetCINamespacedName(ci.Build.String(), b.Play)

	resource := &tkn.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:        namespacedName.Name,
			Namespace:   namespacedName.Namespace,
			Annotations: make(map[string]string),
			Labels:      util.GetCILabels(b.Play),
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
					Name: "flags",
					Type: tkn.ParamTypeArray,
				},
			},
			Steps: b.Steps,
		},
	}
	if b.BuildDocker {
		resource.Spec.Resources.Outputs = []tkn.TaskResource{
			{
				ResourceDeclaration: tkn.ResourceDeclaration{
					Name: ci.ResourceImage,
					Type: tkn.PipelineResourceTypeImage,
				},
			},
		}
		resource.Spec.Params = append(resource.Spec.Params, []tkn.ParamSpec{
			{
				Name: "s3DockerfilePath",
				Type: tkn.ParamTypeString,
			},
			{
				Name: "DOCKERFILE",
				Type: tkn.ParamTypeString,
			},
			{
				Name: "IMAGE",
				Type: tkn.ParamTypeString,
			},
			{
				Name: "CONTEXT",
				Type: tkn.ParamTypeString,
			},
		}...)
	}

	// set the current time in the resource annotations
	resource.Annotations[config.ScheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
	if err := controllerutil.SetControllerReference(b.Play, resource, b.Scheme); err != nil {
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
