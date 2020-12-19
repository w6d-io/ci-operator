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
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// StepSpec defines the desired state of Step
type StepStep struct {
	tkn.Step `json:",inline"`
}

// StepStatus defines the observed state of Step
type StepStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// Step is the Schema for the steps API
type Step struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Step   StepStep   `json:"step"`
	Status StepStatus `json:"status,omitempty"`
}

type Steps []Step

// +kubebuilder:object:root=true

// StepList contains a list of Step
type StepList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           Steps `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Step{}, &StepList{})
}

// Len method for Sort
func (in Steps) Len() int {
	return len(in)
}

// Swap method for Sort
func (in Steps) Swap(i, j int) {
	in[i], in[j] = in[j], in[i]
}

// Less method for Sort
func (in Steps) Less(i, j int) bool {
	var right, left int
	var err error
	right, err = strconv.Atoi(in[j].Annotations[AnnotationOrder])
	if err != nil {
		return false
	}
	left, err = strconv.Atoi(in[i].Annotations[AnnotationOrder])
	if err != nil {
		return true
	}
	return left < right
}
