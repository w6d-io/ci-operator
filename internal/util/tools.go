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
	"context"
	"fmt"
	"net/url"
	ctrl "sigs.k8s.io/controller-runtime"
	"strconv"
	"strings"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"knative.dev/pkg/apis/duck/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Condition returns a kubernetes State
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

// Message returns a kubernetes Message
func Message(c v1beta1.Conditions) string {
	if len(c) == 0 {
		return ""
	}
	return c[0].Message
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

// InNamespace return a client.InNamespace with pipeline namespace
func InNamespace(play *ci.Play) client.InNamespace {
	return client.InNamespace(fmt.Sprintf("p6e-cx-%v", play.Spec.ProjectID))
}

// MatchingLabels return a client.MatchingLabels with pipeline labels
func MatchingLabels(play *ci.Play) client.MatchingLabels {
	return client.MatchingLabels{
		"pipelineid": fmt.Sprintf("%d", play.Spec.PipelineID),
		"projectid":  fmt.Sprintf("%d", play.Spec.ProjectID),
	}
}

// IsPodExist gets a list of pods from namespace and labels and return if there is pod
func IsPodExist(ctx context.Context, r client.Client, play *ci.Play) (bool, error) {
	podList := &corev1.PodList{}

	var opts []client.ListOption
	opts = append(opts, InNamespace(play))
	opts = append(opts, MatchingLabels(play))
	ctrl.Log.WithName("IsPodExist").V(1).Info("DEBUG TEST", "opts", opts)
	if err := r.List(ctx, podList, opts...); err != nil {
		return false, err
	}
	ctrl.Log.V(1).Info("DEBUG TEST", "content", fmt.Sprintf("%+v",
		GetObjectContain(podList)))
	if len(podList.Items) > 0 {
		return true, nil
	}
	return false, nil
}

// GetCINamespacedName return CI namespacedName
func GetCINamespacedName(prefix string, play *ci.Play) types.NamespacedName {
	return types.NamespacedName{
		Name:      fmt.Sprintf("%s-%v-%v", prefix, play.Spec.ProjectID, play.Spec.PipelineID),
		Namespace: fmt.Sprintf("p6e-cx-%v", play.Spec.ProjectID),
	}
}

// GetCINamespacedName2 return CI namespacedName
func GetCINamespacedName2(prefix string, play *ci.Play) types.NamespacedName {
	return types.NamespacedName{
		Name:      fmt.Sprintf("%s-%v", prefix, play.Spec.ProjectID),
		Namespace: fmt.Sprintf("p6e-cx-%v", play.Spec.ProjectID),
	}
}

// GetDeployNamespacedName return CI namespacedName
func GetDeployNamespacedName(prefix string, play *ci.Play) types.NamespacedName {
	return types.NamespacedName{
		Name:      fmt.Sprintf("%s-%v", prefix, play.Spec.ProjectID),
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
	// TODO to remove. DockerURL ahs to be set
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
func GetDockerImageTagRaw(play *ci.Play) (address string, uri string, tag string, err error) {
	var URL *url.URL
	// TODO to remove. DockerURL ahs to be set
	rep := fmt.Sprintf("https://reg-ext.w6d.io/cxcm/%v/%v:%v-%v",
		play.Spec.ProjectID, play.Spec.Name, play.Spec.Commit.SHA[:8], play.Spec.Commit.Ref)
	URL, err = ParseHostURL(rep)
	if err != nil {
		return
	}
	if play.Spec.DockerURL != "" {
		if !strings.HasPrefix(play.Spec.DockerURL, "http") && !strings.Contains(play.Spec.DockerURL, "://") {
			play.Spec.DockerURL = "https://" + play.Spec.DockerURL
		}
		URL, err = ParseHostURL(play.Spec.DockerURL)
		if err != nil {
			return
		}
	}
	partURI := strings.SplitN(URL.Path, ":", 2)
	address = URL.Host
	uri = partURI[0]
	tag = "latest"
	if len(partURI) > 1 && partURI[1] == "" {
		tag = partURI[1]
	}
	return
}

// IgnoreNotExists returns nil on NotExist errors.
// All other values that are not NotExist errors or nil are returned unmodified.
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

// IsBuildStage check whether the stage is build or not
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
	} else {
		parsed, err := url.Parse("http://" + addr)
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

type contextKey int

const (
	contextKeyCorrelationID contextKey = iota
	contextKeyPlay
)

// NewCorrelationIDContext returns a context with correlation id
func NewCorrelationIDContext(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, contextKeyCorrelationID, correlationID)
}

// GetCorrelationIDFromContext gets the value from the context.
func GetCorrelationIDFromContext(ctx context.Context) (string, bool) {
	caller, ok := ctx.Value(contextKeyCorrelationID).(string)
	return caller, ok
}

// NewPlayContext returns a context with correlation id
func NewPlayContext(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, contextKeyPlay, value)
}

// GetPlayFromContext gets the value from the context.
func GetPlayFromContext(ctx context.Context) (string, bool) {
	caller, ok := ctx.Value(contextKeyPlay).(string)
	return caller, ok
}
