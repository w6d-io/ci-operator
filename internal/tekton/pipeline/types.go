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
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

type Pipeline struct {
	Pos            int
	NamespacedName types.NamespacedName
	Labels         map[string]string
	Params         []tkn.ParamSpec
	Tasks          []tkn.PipelineTask
	RunAfter       []string
	Workspaces     []tkn.PipelineWorkspaceDeclaration
	Resources      []tkn.PipelineDeclaredResource
	Play           *ci.Play
	Scheme         *runtime.Scheme
}

type Interface interface {
	SetPipelineUnitTest(logr.Logger) error
	SetPipelineBuild(logr.Logger) error
	SetPipelineDeploy(logr.Logger) error
	SetPipelineIntTest(logr.Logger) error
	SetPipelineClean(logr.Logger) error
	SetPipelineSonar(logr.Logger) error
	SetPipelineE2ETest(logr.Logger) error
}

var _ Interface = &Pipeline{}
