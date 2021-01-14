/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 30/12/2020
*/

package config

import tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"

// Config controller common parameter
type Config struct {
	DefaultDomain string                     `json:"domain"`
	ClusterRole   string                     `json:"cluster_role"`
	PodTemplate   *tkn.PodTemplate           `json:"podTemplate"`
	Workspaces    []tkn.WorkspaceDeclaration `json:"workspaces"`
	Ingress       Ingress                    `json:"ingress"`
	Volume        tkn.WorkspaceBinding       `json:"volume"`
	// Hash is use for provide unpredictable string from an integer
	Hash Hash `json:"hash"`
	// Minio contains all minio information for the connection the could be omitted
	Minio *Minio `json:"minio,omitempty"`
	// DeployPrefix is used to build namespace name where application will be deployed
	// default values is cx
	DeployPrefix string `json:"deploy_prefix"`
}

type Minio struct {
	Host      string `json:"host"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
}

// Ingress struct
type Ingress struct {
	Class  string `json:"class"`
	Prefix string `json:"prefix"`
	Issuer string `json:"issuer"`
}

type Hash struct {
	Salt      string `json:"salt"`
	MinLength int    `json:"min_length"`
}

// config implements Config struct
var (
	config                  = new(Config)
	ScheduledTimeAnnotation = "play.ci.w6d.io/scheduled-at"
)
