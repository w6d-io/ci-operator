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
Created on 16/12/2020
*/

package pipelinerun

import (
	"context"
	"fmt"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/k8s/serviceaccount"
	"github.com/w6d-io/ci-operator/internal/util"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (p *PipelineRun) Parse(log logr.Logger) error {

	for pos, m := range p.Play.Spec.Tasks {
		for name := range m {
			switch name {
			case ci.Build:
				if err := p.SetBuild(pos, log); err != nil {
					return err
				}
			case ci.Clean:
				if err := p.SetClean(pos, log); err != nil {
					return err
				}
			case ci.Deploy:
				if err := p.SetDeploy(pos, log); err != nil {
					return err
				}
			case ci.IntegrationTests:
				if err := p.SetIntTest(pos, log); err != nil {
					return err
				}
			case ci.Sonar:
				if err := p.SetSonar(pos, log); err != nil {
					return err
				}
			case ci.UnitTests:
				if err := p.SetUnitTest(pos, log); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (p *PipelineRun) Create(ctx context.Context, r client.Client, log logr.Logger) error {
	log = log.WithName("Create").WithValues("action", "pipeline-run")
	log.V(1).Info("creating")
	namespacedName := util.GetCINamespacedName("pipeline-run", p.Play)
	pipelineRunResource := &tkn.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:        namespacedName.Name,
			Namespace:   namespacedName.Namespace,
			Annotations: make(map[string]string),
			Labels:      util.GetCILabels(p.Play),
		},
		Spec: tkn.PipelineRunSpec{
			PipelineRef:        &tkn.PipelineRef{Name: util.GetCINamespacedName("pipeline", p.Play).Name},
			ServiceAccountName: util.GetCINamespacedName(serviceaccount.Prefix, p.Play).Name,
			Resources:          p.getPipelineResourceBinding(p.Play),
			Workspaces:         p.getWorkspaceBinding(),
			Params:             p.Params,
			PodTemplate:        p.PodTemplate,
		},
	}
	pipelineRunResource.Annotations[config.ScheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
	if err := controllerutil.SetControllerReference(p.Play, pipelineRunResource, p.Scheme); err != nil {
		return err
	}
	log.V(1).Info(fmt.Sprintf("pipelineRun contains\n%v",
		util.GetObjectContain(pipelineRunResource)))
	if err := r.Create(ctx, pipelineRunResource); err != nil {
		return err
	}
	// All went well
	return nil
}

func (p *PipelineRun) getWorkspaceBinding() []tkn.WorkspaceBinding {
	return []tkn.WorkspaceBinding{
		config.Volume(),
	}
}

func (p *PipelineRun) getPipelineResourceBinding(play *ci.Play) []tkn.PipelineResourceBinding {
	res := []tkn.PipelineResourceBinding{
		{
			Name: ci.ResourceGit,
			ResourceRef: &tkn.PipelineResourceRef{
				Name: util.GetCINamespacedName("pr-git", play).Name,
			},
		},
	}
	if util.IsBuildStage(play) {
		res = append(res, tkn.PipelineResourceBinding{
			Name: ci.ResourceImage,
			ResourceRef: &tkn.PipelineResourceRef{
				Name: util.GetCINamespacedName("pr-image", play).Name,
			},
		})
	}
	return res
}
