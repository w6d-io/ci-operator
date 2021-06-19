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
Created on 07/01/2021
*/

package pipeline

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

const Prefix string = "pipeline"

func (p *Pipeline) Parse(logger logr.Logger) error {
	p.Workspaces = append(p.Workspaces, tkn.PipelineWorkspaceDeclaration{
		Name: config.Volume().Name,
	})
	for pos, m := range p.Play.Spec.Tasks {
		for name := range m {
			switch name {
			case ci.UnitTests:
				_ = p.SetPipelineUnitTest(logger)
			case ci.Build:
				p.Pos = pos
				_ = p.SetPipelineBuild(logger)
			case ci.Deploy:
				_ = p.SetPipelineDeploy(logger)
			case ci.IntegrationTests:
				_ = p.SetPipelineIntTest(logger)
			case ci.Clean:
				_ = p.SetPipelineClean(logger)
			case ci.E2ETests:
				_ = p.SetPipelineE2ETest(logger)
			default:
				_ = p.SetPipelineGeneric(name, logger)
			}
		}
	}
	return nil
}

// Create build the pipeline tekton resource
func (p *Pipeline) Create(ctx context.Context, r client.Client, log logr.Logger) error {
	log = log.WithName("Create").WithValues("action", Prefix)
	log.V(1).Info("creating")

	namespacedName := util.GetCINamespacedName(Prefix, p.Play)
	resources := []tkn.PipelineDeclaredResource{
		{
			Name: ci.ResourceGit,
			Type: tkn.PipelineResourceTypeGit,
		},
	}
	if util.IsBuildStage(p.Play) {
		resources = append(resources, tkn.PipelineDeclaredResource{
			Name: ci.ResourceImage,
			Type: tkn.PipelineResourceTypeImage,
		})
	}
	resource := &tkn.Pipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:        namespacedName.Name,
			Namespace:   namespacedName.Namespace,
			Annotations: make(map[string]string),
			Labels:      util.GetCILabels(p.Play),
		},
		Spec: tkn.PipelineSpec{
			Params:     p.Params,
			Resources:  resources,
			Tasks:      p.Tasks,
			Workspaces: p.Workspaces,
		},
	}

	resource.Annotations[config.ScheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
	if err := controllerutil.SetControllerReference(p.Play, resource, p.Scheme); err != nil {
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

func getWorkspacePipelineTaskBinding() []tkn.WorkspacePipelineTaskBinding {
	var ws []tkn.WorkspacePipelineTaskBinding

	for _, wks := range config.Workspaces() {
		wb := tkn.WorkspacePipelineTaskBinding{
			Name:      wks.Name,
			Workspace: config.Volume().Name,
			SubPath:   wks.Name,
		}
		ws = append(ws, wb)
	}

	return ws
}

// return a pointer on task resource for pipeline
func getPipelineTaskResources(build bool) *tkn.PipelineTaskResources {
	ptr := &tkn.PipelineTaskResources{}

	ptr.Inputs = []tkn.PipelineTaskInputResource{
		{
			Name:     ci.ResourceGit,
			Resource: ci.ResourceGit,
		},
	}
	if build {
		ptr.Outputs = []tkn.PipelineTaskOutputResource{
			{
				Name:     ci.ResourceImage,
				Resource: ci.ResourceImage,
			},
		}
	}
	return ptr
}

func getParamString(name string) tkn.Param {
	return tkn.Param{
		Name: name,
		Value: tkn.ArrayOrString{
			Type:      tkn.ParamTypeString,
			StringVal: "$(params." + name + ")",
		},
	}
}
