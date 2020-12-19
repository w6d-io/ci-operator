/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 22/11/2020
*/

package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sort"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/ghodss/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Ingress struct
type Ingress struct {
	Class  string `json:"class"`
	Prefix string `json:"prefix"`
}

// Config controller common parameter
type Config struct {
	GitImage           string                     `json:"git_image"`
	DefaultDomain      string                     `json:"domain"`
	ServiceAccountName string                     `json:"serviceAccountName,omitempty"`
	PodTemplate        *tkn.PodTemplate           `json:"podTemplate"`
	Workspaces         []tkn.WorkspaceDeclaration `json:"workspaces"`
	Ingress            Ingress                    `json:"ingress"`
	Volume             Volume                     `json:"volume"`
}

// Volume contains the spec for Volume Claim
type Volume struct {
	Size string                            `json:"size"`
	Mode corev1.PersistentVolumeAccessMode `json:"mode"`
}

// New get the filename and fill Config struct
func (c *Config) New(filename string) error {
	log := ctrl.Log.WithName("controllers").WithName("Config")
	log.V(1).Info("read config file")
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error(err, "error reading the configuration")
		return err
	}
	if err := yaml.Unmarshal(data, c); err != nil {
		log.Error(err, "Error unmarshal the configuration")
	}
	return nil
}

// Validate returns an error if is mandatory Config missing
func (c *Config) Validate() error {
	var unset []string
	for _, f := range []struct {
		v, name string
	}{
		{c.DefaultDomain, "domain"},
		{c.Ingress.Class, "ingress.class"},
	} {
		if f.v == "" {
			unset = append(unset, f.name)
		}
	}
	if len(unset) > 0 {
		sort.Strings(unset)
		return fmt.Errorf("found unset image flags: %s", unset)
	}
	return nil
}

func New() *WFType {
	return &WFType{
		Index:   0,
		Creates: []CIFunc{},
	}
}

func (wf *WFType) Add(ciFunc CIFunc) error {
	wf.Creates = append(wf.Creates, ciFunc)
	return nil
}

// Cfg implements Config struct
var Cfg = new(Config)

var (
	scheduledTimeAnnotation = "play.ci.w6d.io/scheduled-at"
)

// WFInterface implements all Workflow methods
type WFInterface interface {
	// Add func in Creates
	Add(ciFunc CIFunc) error

	// SetGitCreate create the git pipeline resource type
	SetGitCreate(p ci.Play, r *PlayReconciler) error
	// SetGitCreate create the image pipeline resource type
	SetImageCreate(p ci.Play, r *PlayReconciler) error

	// SetTask execute the Task according the TaskType
	SetTask(ctx context.Context, p ci.Play, r *PlayReconciler) error
	// Build  implements the build Tekton task
	Build(ctx context.Context, p ci.Play, r *PlayReconciler) error
	// UnitTest  implements the unit test Tekton task
	UnitTest(ctx context.Context, p ci.Play, r *PlayReconciler) error
	// IntTest implements the integration test Tekton task
	IntTest(ctx context.Context, p ci.Play, r *PlayReconciler) error
	// Deploy implements the deploy Tekton task
	Deploy(ctx context.Context, p ci.Play, r *PlayReconciler) error
	// Sonar implements the sonar Tekton task
	Sonar(ctx context.Context, p ci.Play, r *PlayReconciler) error
	// Clean implements the clean Tekton task
	Clean(ctx context.Context, p ci.Play, r *PlayReconciler) error

	SetPipeline(play ci.Play, r *PlayReconciler) error
	SetPipelineRun(play ci.Play, r *PlayReconciler) error

	Run(ctx context.Context, r client.Client, log logr.Logger) error
}

// WFType contains all tekton resource to create
type WFType struct {
	Index   int
	Creates []CIFunc
}

// CIInterface implements the CI method to create tekton resource
type CIInterface interface {
	// Create will create tekton resource
	Create(ctx context.Context, r client.Client, logger logr.Logger) error
}

// CIFunc is a function that implements the CIInterface
type CIFunc func(ctx context.Context, c client.Client, logger logr.Logger) error

type Pipeline struct {
	OwnerReferences []metav1.OwnerReference
	NamespacedName  types.NamespacedName
	Labels          map[string]string
	Params          []tkn.ParamSpec
	Tasks           []tkn.PipelineTask
	RunAfter        []string
	Workspaces      []tkn.PipelineWorkspaceDeclaration
	Resources       []tkn.PipelineDeclaredResource
}

type PipelineInterface interface {
	SetPipelineUnitTest(play ci.Play, r *PlayReconciler) error
	SetPipelineBuild(play ci.Play, r *PlayReconciler) error
	SetPipelineDeploy(play ci.Play, r *PlayReconciler) error
	SetPipelineIntTest(play ci.Play, r *PlayReconciler) error
	SetPipelineClean(play ci.Play, r *PlayReconciler) error
	SetPipelineSonar(play ci.Play, r *PlayReconciler) error
}
