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
Created on 22/11/2020
*/

package play

import (
	"context"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"

	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal/k8s/rbac"
	"github.com/w6d-io/ci-operator/internal/k8s/sa"
	"github.com/w6d-io/ci-operator/internal/tekton/pipeline"
	"github.com/w6d-io/ci-operator/internal/tekton/pipelineresource"
	"github.com/w6d-io/ci-operator/internal/tekton/pipelinerun"
	"github.com/w6d-io/ci-operator/internal/tekton/task"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// WFInterface implements all Workflow methods
type WFInterface interface {
	// Add func in Creates list
	Add(ciFunc CIFunc) error

	CreateValues(context.Context, *ci.Play, logr.Logger) error

	// ServiceAccount creates the serviceAccount resource use for running pipeline and handle secret
	ServiceAccount(*ci.Play, logr.Logger) error
	// Rbac creates the RoleBinding resource bind with the ServiceAccount
	Rbac(*ci.Play, logr.Logger) error
	// GitSecret create the secret resource for cloning
	GitSecret(*ci.Play, logr.Logger) error
	// DockerCredSecret creates the docker secret resource for pull and push image
	DockerCredSecret(*ci.Play, logr.Logger) error
	// MinIOSecret creates the s3cfg file in a secret resource
	MinIOSecret(*ci.Play, logr.Logger) error

	// KubeConfigSecret creates the ~/.kube/config file in a secret resource
	KubeConfigSecret(*ci.Play, logr.Logger) error

	// VaultSecret creates the ~/.kube/config file in a secret resource
	//VaultSecret(*ci.Play, logr.Logger) error

	// SetGitCreate creates the git pipeline resource type
	SetGitCreate(*ci.Play, logr.Logger) error
	// SetImageCreate creates the image pipeline resource type
	SetImageCreate(*ci.Play, logr.Logger) error

	// SetTask executes the Task according the TaskType
	SetTask(context.Context, *ci.Play, logr.Logger) error
	// Build  implements the build Tekton task
	//Build(context.Context, *ci.Play, logr.Logger) error
	//// UnitTest  implements the unit test Tekton task
	//UnitTest(context.Context, *ci.Play, logr.Logger) error
	//// IntTest implements the integration test Tekton task
	//IntTest(context.Context, *ci.Play, logr.Logger) error
	//// Deploy implements the deploy Tekton task
	//Deploy(context.Context, *ci.Play, logr.Logger) error
	//// Sonar implements the sonar Tekton task
	//Sonar(context.Context, *ci.Play, logr.Logger) error
	//// Clean implements the clean Tekton task
	//Clean(context.Context, *ci.Play, logr.Logger) error

	// SetPipeline implements the pipeline Tekton resource
	SetPipeline(*ci.Play, logr.Logger) error
	// SetPipelineRun implements the pipeline run Tekton resource
	SetPipelineRun(*ci.Play, logr.Logger) error

	// Run executes all create method in the list
	Run(context.Context, client.Client, logr.Logger) error
}

var _ WFInterface = &WFType{}

// WFType contains all tekton resource to create
type WFType struct {
	// Creates is the list of create methods
	Creates []CIFunc
	Client  client.Client
	Scheme  *runtime.Scheme
	Params  map[string][]ci.ParamSpec
}

// CIInterface implements the CI method to create tekton resource
type CIInterface interface {
	// Create will create tekton resource
	Create(ctx context.Context, r client.Client, logger logr.Logger) error
}

var _ CIInterface = &pipelineresource.ImagePR{}
var _ CIInterface = &pipelineresource.GitPR{}
var _ CIInterface = &sa.CI{}
var _ CIInterface = &sa.Deploy{}
var _ CIInterface = &rbac.CI{}
var _ CIInterface = &rbac.Deploy{}
var _ CIInterface = &task.BuildTask{}
var _ CIInterface = &task.CleanTask{}
var _ CIInterface = &task.DeployTask{}
var _ CIInterface = &task.IntTestTask{}
var _ CIInterface = &task.UnitTestTask{}
var _ CIInterface = &pipeline.Pipeline{}
var _ CIInterface = &pipelinerun.PipelineRun{}

// CIFunc is a function that implements the CIInterface
type CIFunc func(ctx context.Context, c client.Client, logger logr.Logger) error
