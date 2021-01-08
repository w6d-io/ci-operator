/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
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
	SetPipelineUnitTest(*ci.Play, logr.Logger) error
	SetPipelineBuild(*ci.Play, logr.Logger) error
	SetPipelineDeploy(*ci.Play, logr.Logger) error
	SetPipelineIntTest(*ci.Play, logr.Logger) error
	SetPipelineClean(*ci.Play, logr.Logger) error
	SetPipelineSonar(*ci.Play, logr.Logger) error
}

var _ Interface = &Pipeline{}
