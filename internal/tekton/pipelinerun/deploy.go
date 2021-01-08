/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 17/12/2020
*/

package pipelinerun

import (
	"fmt"
	"github.com/go-logr/logr"
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/k8s/secrets"
	"github.com/w6d-io/ci-operator/internal/util"
	corev1 "k8s.io/api/core/v1"
)

// SetDeploy builds the params adn add the minio volume for the tekton pipelineRun resource
// param pos : use to get task part from Play to get variables
func (p *PipelineRun) SetDeploy(pos int, log logr.Logger) error {
	log = log.WithName("SetDeploy").WithValues("action", "pipeline-run")
	log.V(1).Info("set deploy pipeline run params")

	task := p.Play.Spec.Tasks[pos][ci.Deploy]
	var flags []string
	if len(task.Variables) != 0 {
		for key, val := range task.Variables {
			flags = append(flags, "--set")
			flags = append(flags, key+"="+val)
		}
	}
	p.Params = append(p.Params, tkn.Param{
		Name: "deploy_flags",
		Value: tkn.ArrayOrString{
			Type:     tkn.ParamTypeArray,
			ArrayVal: flags,
		},
	}, tkn.Param{
		Name: "deploy_s3valuepath",
		Value: tkn.ArrayOrString{
			Type: tkn.ParamTypeString,
			StringVal: fmt.Sprintf("%v/%v/values.yaml",
				p.Play.Spec.ProjectID,
				p.Play.Spec.PipelineID),
		},
	}, tkn.Param{
		Name: "deploy_values",
		Value: tkn.ArrayOrString{
			Type: tkn.ParamTypeString,
			StringVal: fmt.Sprintf("%s/%s/values.yaml", config.GetMinio().GetBucket(),
				config.GetWorkspacePath("values")),
		},
	}, tkn.Param{
		Name: "deploy_namespace",
		Value: tkn.ArrayOrString{
			Type: tkn.ParamTypeString,
			// TODO put the prefix in config
			StringVal: util.GetDeployNamespacedName("cx", p.Play).Namespace,
		},
	}, tkn.Param{
		Name: "deploy_release_name",
		Value: tkn.ArrayOrString{
			Type:      tkn.ParamTypeString,
			StringVal: util.GetDeployNamespacedName("cx", p.Play).Namespace,
		},
	})

	p.PodTemplate = &tkn.PodTemplate{}
	if config.GetMinio().Host != "" {
		p.PodTemplate.Volumes = []corev1.Volume{
			{
				Name: secrets.MinIOPrefixSecret,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: util.GetCINamespacedName(secrets.MinIOPrefixSecret, p.Play).Name,
					},
				},
			},
		}
	}

	return nil
}
