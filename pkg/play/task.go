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
Created on 23/11/2020
*/

package play

import (
	"context"
	"github.com/go-logr/logr"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal/tekton/task"
)

// SetTask executes the Task according the TaskType
//   - Build     => create a build Tekton Task
//   - Clean     => create a cleaning Tekton Task
//   - Deploy    => create a deployment Tekton Task
//   - IntTest   => create a integration test Tekton Task
//   - Sonar     => create a sonar Tekton Task
//   - UnitTests => create a unit test Tekton Task
//   - E2ETests  => create a e2e test Tekton Task
func (wf *WFType) SetTask(ctx context.Context, p *ci.Play, logger logr.Logger) error {
	log := logger.WithName("SetTask")
	log.Info("Build tasks")
	var t = task.Task{
		Client: wf.Client,
		Play:   p,
		Scheme: wf.Scheme,
	}
	if err := t.Parse(ctx, logger); err != nil {
		return err
	}
	for key := range t.Params {
		wf.Params[key] = t.Params[key]
	}
	for _, create := range t.Creates {
		if err := wf.Add(create); err != nil {
			return err
		}
	}
	return nil
}
