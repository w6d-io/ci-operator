/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 21/11/2020
*/

package controllers

import (
	"bytes"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
	"net/url"
	"strconv"
	"strings"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/pkg/apis/duck/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Condition returns a Play State
func Condition(c v1beta1.Conditions) (status ci.State) {
	if len(c) == 0 {
		return "---"
	}

	switch c[0].Status {
	case corev1.ConditionFalse:
		status = ci.Failed
	case corev1.ConditionTrue:
		status = ci.Succeeded
	case corev1.ConditionUnknown:
		status = ci.Running
	}
	if c[0].Reason != "" {
		if c[0].Reason == "PipelineRunCancelled" || c[0].Reason == "TaskRunCancelled" {
			status = ci.Cancelled
		}
	}
	return
}

// IsPipelineRunning return whether or not the pipeline is running
func IsPipelineRunning(pr tkn.PipelineRun) bool {

	nonRunningState := map[ci.State]bool{
		ci.Failed:    true,
		ci.Cancelled: true,
		ci.Succeeded: true,
	}
	if _, ok := nonRunningState[Condition(pr.Status.Conditions)]; ok {
		return false
	}
	return true
}

func InNamespace(play ci.Play) client.InNamespace {
	return client.InNamespace(fmt.Sprintf("p6e-cx-%v", play.Spec.ProjectID))
}

// CxCINamespacedName return CI namespacedName
func CxCINamespacedName(prefix string, play ci.Play) types.NamespacedName {
	return types.NamespacedName{
		Name:      fmt.Sprintf("%s-%v-%v", prefix, play.Spec.ProjectID, play.Spec.PipelineID),
		Namespace: fmt.Sprintf("p6e-cx-%v", play.Spec.ProjectID),
	}
}

// CxCINamespacedName return CI namespacedName
func CxDeployNamespacedName(prefix string, play ci.Play) types.NamespacedName {
	return types.NamespacedName{
		Name:      fmt.Sprintf("%s-%v-%v", prefix, play.Spec.ProjectID, play.Spec.PipelineID),
		Namespace: fmt.Sprintf("%s-%v-%v", prefix, play.Spec.Environment, play.Spec.ProjectID),
	}
}

// CxCILabels return CI label
func CxCILabels(p ci.Play) map[string]string {
	return map[string]string{
		// TODO put label key in configuration
		"projectid":  strconv.Itoa(int(p.Spec.ProjectID)),
		"pipelineid": strconv.Itoa(int(p.Spec.PipelineID)),
	}
}

// CxDockerImageName return the URL of the Docker repository
func CxDockerImageName(play ci.Play) (*url.URL, error) {
	// TODO find a way to get the user docker registry
	URL, err := url.Parse(strings.ToLower(fmt.Sprintf("reg-ext.w6d.io/cxcm/%v/%v:%v-%v",
		play.Spec.ProjectID, play.Spec.Name, play.Spec.Commit.SHA[:8], play.Spec.Commit.Ref)))
	if err != nil {
		return nil, err
	}
	return URL, nil
}

// IgnoreNotFound returns nil on NotFound errors.
// All other values that are not NotFound errors or nil are returned unmodified.
func IgnoreNotExists(err error) error {
	if err == nil ||
		(strings.HasPrefix(err.Error(), "Index with name field:") &&
			strings.HasSuffix(err.Error(), "does not exist")) {
		return nil
	}
	return err
}

// CIOwnerReferences return a CI owner reference
func CIOwnerReference(play ci.Play) metav1.OwnerReference {
	ownerRef := metav1.OwnerReference{}
	ownerRef.Name = play.Name
	ownerRef.Kind = play.Kind
	ownerRef.APIVersion = play.APIVersion
	ownerRef.UID = play.UID

	return ownerRef
}

// GetObjectContain ...
func GetObjectContain(obj runtime.Object) string {
	s := json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, json.SerializerOptions{Yaml: true})
	buf := new(bytes.Buffer)
	s.Encode(obj, buf)
	return buf.String()
}
