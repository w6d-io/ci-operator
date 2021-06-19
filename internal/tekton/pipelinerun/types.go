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
	"github.com/go-logr/logr"
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PipelineRun struct {
	internal.WorkFlowStruct

	// Params contains list of param for PipelineRun tekton resource
	Params []tkn.Param

	// PodTemplate is use for
	// - volumes
	// - nodeSelector
	// - tolerations
	PodTemplate *tkn.PodTemplate

	GenericParams map[string][]ci.ParamSpec
}

const (
	Prefix string = "pipeline-run"
)

type Interface interface {
	SetBuild(int, logr.Logger) error
	SetClean(int, logr.Logger) error
	SetDeploy(int, logr.Logger) error
	SetE2ETest(int, logr.Logger) error
	SetGeneric(int, ci.TaskType, logr.Logger) error
	SetIntTest(int, logr.Logger) error
	SetUnitTest(int, logr.Logger) error

	Parse(logr.Logger) error
	GetNamespace(ci.Task) string
	Create(context.Context, client.Client, logr.Logger) error
}

var _ Interface = &PipelineRun{}
