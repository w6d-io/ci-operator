/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PlayGroupKind is group kind used in validator
var (
	PlayGroupKind = schema.GroupKind{Group: "ci.w6d.io", Kind: "Play"}
)


// PlaySpec defines the desired state of Play
type PlaySpec struct {
	// Name of project
	Name string `json:"name"`

	// Stack of the project
	Stack Stack `json:"stack"`

	// Scope contains the name of scope and list of projects id
	Scope Scope `json:"scope"`

	// Environment contains application environment
	Environment string `json:"environment"`

	// ProjectID contains the project ID
	ProjectID int64 `json:"project_id"`

	// PipelineID contains the ID of pipeline for the project
	PipelineID int64 `json:"pipeline_id"`

	// RepoURL contains the git repository url
	RepoURL string `json:"repo_url"`

	// Token contains the token for git clone
	// +optional
	Token string `json:"token,omitempty"`

	// Commit contains all git information
	Commit Commit `json:"commit"`

	// Sonar contains the sonarqube token
	// +optional
	Sonar string `json:"sonar,omitempty"`

	// Domain contains the url for exposition
	// +optional
	// +kubebuilder:validation:Pattern=^([A-Za-z0–9-]+\.)+[A-Za-z][A-Za-z]+$
	Domain string `json:"domain,omitempty"`

	// Tasks contains the list of task to be created by Play
	Tasks []map[TaskType]Task `json:"tasks"`

	// Dependencies contains a list of Dependency ie: MongoDb or PostgreSQL
	// +optional
	Dependencies map[DependencyType]Dependency `json:"dependencies,omitempty"`

	// Vmx toggle for use vmx node
	// +optional
	Vmx bool `json:"vmx,omitempty"`
}

// PlayStatus defines the observed state of Play
type PlayStatus struct {
	// PipelineRunName contains the pipeline run name created by play
	// +optional
	PipelineRunName string `json:"pipeline_run_name,omitempty"`

	// State contains the current state of this Play resource.
	// States Running, Failed, Succeeded, Errored
	// +optional
	State State `json:"state,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.state"
// +kubebuilder:printcolumn:name="PipelineRun",type="string",priority=1,JSONPath=".status.pipeline_run"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="CreationTimestamp is a timestamp representing the server time when this object was created. It is not guaranteed to be set in happens-before order across separate operations. Clients may not set this value. It is represented in RFC3339 form and is in UTC."
// +kubebuilder:subresource:status
// Play is the Schema for the plays API
type Play struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PlaySpec   `json:"spec"`
	Status PlayStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PlayList contains a list of Play
type PlayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Play `json:"items"`
}

// Commit contains all git information
type Commit struct {
	// SHA contains git commit SHA
	SHA string `json:"sha"`
	// Ref contains git commit reference
	Ref string `json:"ref"`
	// Message contains commit message
	// +optional
	Message string `json:"message,omitempty"`
}

// Task is what actions and/or configuration the task can be contains
type Task struct {
	// Script is a list of command to execute in the task
	// +optional
	Script Script `json:"script,omitempty"`

	// Env is the map of environment variable for the task
	// +optional
	Variables fields.Set `json:"variables,omitempty"`

	// Docker contains information for docker build
	// +optional
	Docker Docker `json:"docker,omitempty"`
}

// Dependency struct contains env and service for the dependencies
type Dependency struct {
	// Env is environmental variable for this dependency
	Env fields.Set `json:"env"`

	// Services contain a list of host and port to expose
	// +optional
	Services []Service `json:"services,omitempty"`
}

// NameValue struct for env type format kubernetes format
type NameValue struct {
	Name   string `json:"name"`
	Values string `json:"values"`
}

// Docker structure contains information for docker build
type Docker struct {
	// Filepath contains the dockerfile full path
	// +optional
	Filepath string `json:"filepath,omitempty"`
	// Context contains the docker build context
	// +optional
	Context string `json:"context,omitempty"`
}

// TaskType is the list of task granted
// +kubebuilder:validation:Enum=tests-unit;build;tests-integration;deploy;clean;sonar
type TaskType string

const (
	// UnitTests is the task type for unit tests"
	UnitTests TaskType = "unit-tests"

	// Sonar is the task type for Sonar scan"
	Sonar TaskType = "sonar"

	// Build is the task type for build"
	Build TaskType = "build"

	// IntegrationTests is the task type for integration tests"
	IntegrationTests TaskType = "integration-tests"

	// Deploy is the task type for deploy"
	Deploy TaskType = "deploy"

	// Clean is the task type for clean"
	Clean TaskType = "clean"
)

// Scope is use for gathering project
type Scope struct {
	// Name of the scope
	Name string `json:"name"`

	// Projects is the list of project id in this scope
	Projects int64 `json:"projects"`
}

// Stack contains the language and package of the source
type Stack struct {
	// Language contains the repository language
	Language string `json:"language"`

	// Package contains the package use in application
	// +optional
	Package string `json:"package"`
}

func (in Stack) String() string {
	return in.Language + "/" + in.Package
}

// ServiceElem is the contain of service
// +kubebuilder:validation:Enum=Host;Port
type ServiceElement string

// Service struct for dependencies
type Service map[ServiceElement]string

// DependencyType contain list of dependencies managed
// +kubebuilder::validation:Enum=mongodb;postgresql
type DependencyType string

// Script type
type Script []string

// State type
type State string

const (
	// Creating means that tekton resource creation is in progress
	Creating State = "Creating"

	// Queued means that the PipelineRun not applied yet due to limitation
	Queued State = "Queued"

	// Running signifies at least on Step of the Task is running
	Running State = "Running"

	// Failed signifies at least on Step of the Task is failed
	Failed State = "Failed"

	// Succeeded signifies that all Task is success
	Succeeded State = "Succeeded"

	// Cancelled signifies that a TaskRun or PipelineRun has been cancelled
	Cancelled State = "Cancelled"

	// Errored signifies that at least one tekton resource couldn't be created
	Errored State = "Errored"
)

func init() {
	SchemeBuilder.Register(&Play{}, &PlayList{})
}

func (t TaskType) String() string {
	return string(t)
}

// File mode
const (
	// FileMode0755 int32 = 0755
	// FileMode0644 int32 = 0644
	FileMode0444 int32 = 0444
	// FileMode0400 int32 = 0400
)

// Annotations keys
const (
	AnnotationOrder    string = "ci.w6d.io/order"
	AnnotationTask     string = "ci.w6d.io/task"
	AnnotationLanguage string = "ci.w6d.io/language"
	AnnotationPackage  string = "ci.w6d.io/package"
)

// ResourceNames
const (
	ResourceGit   string = "source"
	ResourceImage string = "builddocker"
)
