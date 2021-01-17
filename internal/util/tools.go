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
Created on 21/11/2020
*/

package util

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

func InNamespace(play *ci.Play) client.InNamespace {
	return client.InNamespace(fmt.Sprintf("p6e-cx-%v", play.Spec.ProjectID))
}

// GetCINamespacedName return CI namespacedName
func GetCINamespacedName(prefix string, play *ci.Play) types.NamespacedName {
	return types.NamespacedName{
		Name:      fmt.Sprintf("%s-%v-%v", prefix, play.Spec.ProjectID, play.Spec.PipelineID),
		Namespace: fmt.Sprintf("p6e-cx-%v", play.Spec.ProjectID),
	}
}

// GetCINamespacedName return CI namespacedName
func GetDeployNamespacedName(prefix string, play *ci.Play) types.NamespacedName {
	return types.NamespacedName{
		Name:      fmt.Sprintf("%s-%v-%v", prefix, play.Spec.ProjectID, play.Spec.PipelineID),
		Namespace: fmt.Sprintf("%s-%v-%v", prefix, play.Spec.Environment, play.Spec.ProjectID),
	}
}

// GetCILabels return CI label
func GetCILabels(p *ci.Play) map[string]string {
	return map[string]string{
		// TODO put label key in configuration
		"projectid":  strconv.Itoa(int(p.Spec.ProjectID)),
		"pipelineid": strconv.Itoa(int(p.Spec.PipelineID)),
	}
}

// GetDockerImageTag return the URL of the Docker repository
func GetDockerImageTag(play *ci.Play) (*url.URL, error) {
	// TODO find a way to get the user docker registry
	rawURL := fmt.Sprintf("reg-ext.w6d.io/cxcm/%v/%v:%v-%v",
		play.Spec.ProjectID, play.Spec.Name, play.Spec.Commit.SHA[:8], play.Spec.Commit.Ref)
	if play.Spec.DockerURL != "" {
		rawURL = play.Spec.DockerURL
	}
	URL, err := url.Parse(strings.ToLower(rawURL))
	if err != nil {
		return nil, err
	}
	return URL, nil
}

// GetDockerImageTagRaw return the Docker repository
func GetDockerImageTagRaw(play *ci.Play) (address string, tag string, err error) {

	rep := fmt.Sprintf("https://reg-ext.w6d.io/cxcm/%v/%v:%v-%v",
		play.Spec.ProjectID, play.Spec.Name, play.Spec.Commit.SHA[:8], play.Spec.Commit.Ref)
	URL, err := ParseHostURL(rep)
	if err != nil {
		return "", "", err
	}
	if play.Spec.DockerURL != "" {
		if !strings.HasPrefix(play.Spec.DockerURL, "http") {
			play.Spec.DockerURL = "https://" + play.Spec.DockerURL
		}
		URL, err = ParseHostURL(play.Spec.DockerURL)
		if err != nil {
			return
		}
	}

	partAddress := strings.SplitN(URL.Path, ":", 2)
	address = partAddress[0]
	tag = partAddress[1]
	if partAddress[1] == "" {
		tag = "latest"
	}
	return
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

// GetObjectContain ...
func GetObjectContain(obj runtime.Object) string {
	s := json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, json.SerializerOptions{Yaml: true})
	buf := new(bytes.Buffer)
	if err := s.Encode(obj, buf); err != nil {
		return "<ERROR>\n"
	}
	return buf.String()
}

// GetStage build the tasks
func IsBuildStage(play *ci.Play) bool {
	if strings.ToLower(play.Spec.Stack.Language) == "android" ||
		strings.ToLower(play.Spec.Stack.Language) == "ios" {
		return false
	}
	for _, t := range play.Spec.Tasks {
		for taskType := range t {
			if taskType == ci.Build {
				return true
			}
		}
	}
	return false
}


// ParseHostURL parses a url string, validates the string is a host url, and
// returns the parsed URL
func ParseHostURL(host string) (*url.URL, error) {
	protoAddrParts := strings.SplitN(host, "://", 2)
	if len(protoAddrParts) == 1 {
		return nil, fmt.Errorf("unable to parse docker host `%s`", host)
	}

	var basePath string
	proto, addr := protoAddrParts[0], protoAddrParts[1]
	if proto == "tcp" {
		parsed, err := url.Parse("tcp://" + addr)
		if err != nil {
			return nil, err
		}
		addr = parsed.Host
		basePath = parsed.Path
	}
	return &url.URL{
		Scheme: proto,
		Host:   addr,
		Path:   basePath,
	}, nil
}