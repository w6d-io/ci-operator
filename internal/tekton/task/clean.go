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
	"github.com/go-logr/logr"
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

// CleanTask task struct for CI
type CleanTask struct {
	Meta
}

// Clean create the clean Tekton Task resource
func (t *Task) Clean(ctx context.Context, logger logr.Logger) error {
	log := logger.WithName("Clean").WithValues("task", ci.Clean)
	// get the task
	log.V(1).Info("get task")
	// TODO get task from index in WFType
	s := &Step{
		Index:    t.Index,
		PlaySpec: t.Play.Spec,
		Client:   t.Client,
		TaskType: ci.Clean,
	}
	steps, err := s.GetSteps(ctx, logger)
	if err != nil {
		log.Error(err, "get steps failed")
		return err
	}
	if len(steps) == 0 {
		return fmt.Errorf("no step found for %s", ci.Clean)
	}
	clean := &CleanTask{
		Meta: Meta{
			Steps:  steps,
			Play:   t.Play,
			Scheme: t.Scheme,
		},
	}

	log.V(1).Info("add create in workflow")
	if err := t.Add(clean.Create); err != nil {
		return err
	}
	return nil
}

func (u *CleanTask) Create(ctx context.Context, r client.Client, log logr.Logger) error {
	log = log.WithName("Create").WithValues("task", ci.Clean)
	log.V(1).Info("creating")
	// build Tekton Task resource
	namespacedName := util.GetCINamespacedName(ci.Clean.String(), u.Play)
	taskResource := &tkn.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:        namespacedName.Name,
			Namespace:   namespacedName.Namespace,
			Annotations: make(map[string]string),
			Labels:      util.GetCILabels(u.Play),
		},
		Spec: tkn.TaskSpec{
			Workspaces: config.Workspaces(),
			Params: []tkn.ParamSpec{
				{
					Name: "namespace",
					Type: tkn.ParamTypeString,
				},
				{
					Name: "release_name",
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
	log.V(1).Info(fmt.Sprintf("task contains\n%v", util.GetObjectContain(taskResource)))
	if err := r.Create(ctx, taskResource); err != nil {
		return err
	}
	// All went well
	return nil
}
