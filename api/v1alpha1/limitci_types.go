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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LimitCiSpec defines the desired state of LimitCi
type LimitCiSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Concurrent contains the number of pipeline running in concurrency
	Concurrent int64 `json:"concurrent,omitempty"`
}

// LimitCiStatus defines the observed state of LimitCi
type LimitCiStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LimitCi is the Schema for the limitcis API
type LimitCi struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LimitCiSpec   `json:"spec,omitempty"`
	Status LimitCiStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LimitCiList contains a list of LimitCi
type LimitCiList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LimitCi `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LimitCi{}, &LimitCiList{})
}
