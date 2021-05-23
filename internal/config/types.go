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
Created on 30/12/2020
*/

package config

import (
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/w6d-io/ci-operator/pkg/webhook"
	corev1 "k8s.io/api/core/v1"
)

// Config controller common parameter
type Config struct {
	DefaultDomain string                     `json:"domain" yaml:"domain"`
	ClusterRole   string                     `json:"cluster_role" yaml:"cluster_role"`
	PodTemplate   *tkn.PodTemplate           `json:"podTemplate" yaml:"podTemplate"`
	Workspaces    []tkn.WorkspaceDeclaration `json:"workspaces" yaml:"workspaces"`
	Ingress       Ingress                    `json:"ingress" yaml:"ingress"`
	Volume        tkn.WorkspaceBinding       `json:"volume" yaml:"volume"`
	Namespace     string                     `json:"namespace" yaml:"namespace"`
	ValuesRef     Values                     `json:"values,omitempty" yaml:"values,omitEmpty"`

	// Hash is use for provide unpredictable string from an integer
	Hash *Hash `json:"hash" yaml:"hash"`

	// Minio contains all minio information for the connection the could be omitted
	Minio *Minio `json:"minio,omitempty" yaml:"minio,omitempty"`

	// DeployPrefix is used to build namespace name where application will be deployed
	// default values is cx
	DeployPrefix string `json:"deploy_prefix" yaml:"deploy_prefix"`

	// WebHooks contains a list of WebHook where payload will be send
	Webhooks []webhook.Webhook `json:"webhooks" yaml:"webhooks"`

	// Vault address
	Vault *Vault `json:"vault,omitempty"`

	EnvPrefix string `json:"envPrefix,omitempty" yaml:"envPrefix,omitempty"`
}

type Values struct {
	// DeployRef contains the configmap information for templating
	DeployRef *corev1.ConfigMapKeySelector `json:"deploy,omitempty"`
}

type Minio struct {
	Host      string `json:"host" yaml:"host"`
	AccessKey string `json:"access_key" yaml:"access_key"`
	SecretKey string `json:"secret_key" yaml:"secret_key"`
	Bucket    string `json:"bucket" yaml:"bucket"`
}

// Ingress struct
type Ingress struct {
	Class  string `json:"class" yaml:"class"`
	Prefix string `json:"prefix" yaml:"prefix"`
	Issuer string `json:"issuer" yaml:"issuer"`
}

type Hash struct {
	Salt      string `json:"salt" yaml:"salt"`
	MinLength int    `json:"min_length" yaml:"min_length"`
}

type Vault struct {
	Host  string `json:"host" yaml:"host"`
	Token string `json:"token" yaml:"token"`
}

// config implements Config struct
var (
	config                  = new(Config)
	ScheduledTimeAnnotation = "play.ci.w6d.io/scheduled-at"
)
