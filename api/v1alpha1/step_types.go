/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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
type StepSpec struct {
	tkn.Step `json:",inline"`
}

// StepStatus defines the observed state of Step
type StepStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// ParamSpec defines arbitrary parameters needed beyond typed inputs (such as
// resources). Parameter values are provided by users as inputs on a TaskRun
// or PipelineRun. It contains also the values template of the params
type ParamSpec struct {

	tkn.ParamSpec `json:",inline"`

	Value string  `json:"value,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Step is the Schema for the steps API
type Step struct {
	metav1.TypeMeta        `json:",inline"`
	metav1.ObjectMeta      `json:"metadata,omitempty"`

	// Parameters declares parameters passed to this task.
	// +optional
	Params []ParamSpec `json:"params,omitempty"`

	Step   StepSpec        `json:"step,omitempty"`
	Status StepStatus      `json:"status,omitempty"`
}

type Steps []Step

//+kubebuilder:object:root=true

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
