/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 16/12/2020
*/

package pipelinerun

import (
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/url"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

type PipelineRun struct {
	Params []tkn.Param

	GetObjectContain     func(runtime.Object) string
	GetNamespacedName    func(string, ci.Play) types.NamespacedName
	DeployNamespacedName func(string, ci.Play) types.NamespacedName
	GetLabels            func(ci.Play) map[string]string
	GetOwnerReferences   func(ci.Play) metav1.OwnerReference
	IsBuildStage         func(ci.Play) bool
	GetDockerImage       func(ci.Play) (*url.URL, error)
	GetWorkspace         func(string, []tkn.WorkspaceDeclaration) string

	ServiceAccount          string
	scheduledTimeAnnotation string
	Play                    ci.Play
	Volume                  Volume
	WorkspacesDeclaration   []tkn.WorkspaceDeclaration
}

// Volume contains the spec for Volume Claim
type Volume struct {
	Size string                            `json:"size"`
	Mode corev1.PersistentVolumeAccessMode `json:"mode"`
}
