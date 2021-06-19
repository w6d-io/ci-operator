/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 21/05/2021
*/

package task

import (
	"fmt"
	"strconv"

	"github.com/w6d-io/ci-operator/internal/config"

	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// BuildAndGetPredefinedEnv build the variables environment for all steps
func BuildAndGetPredefinedEnv(play *ci.Play) (envVars []corev1.EnvVar) {
	envVars = append(envVars, BuildCommonPredefinedEnv(play)...)
	envVars = append(envVars, BuildTaskPredefinedEnv(play.Spec.Tasks)...)
	envVars = append(envVars, BuildConfigPredefinedEnv()...)
	return
}

// BuildCommonPredefinedEnv return the basics commons environment variable
func BuildCommonPredefinedEnv(play *ci.Play) (envVars []corev1.EnvVar) {
	envVars = append(envVars, corev1.EnvVar{
		Name:  config.GetEnvPrefix("project") + "NAME",
		Value: play.Spec.Name,
	})
	envVars = append(envVars, corev1.EnvVar{
		Name:  config.GetEnvPrefix("project") + "LANGUAGE",
		Value: play.Spec.Stack.Language,
	})
	envVars = append(envVars, corev1.EnvVar{
		Name:  config.GetEnvPrefix("project") + "PACKAGE",
		Value: play.Spec.Stack.Package,
	})
	envVars = append(envVars, corev1.EnvVar{
		Name:  config.GetEnvPrefix() + "ENVIRONMENT",
		Value: play.Spec.Environment,
	})
	envVars = append(envVars, corev1.EnvVar{
		Name:  config.GetEnvPrefix() + "PROJECT_ID",
		Value: fmt.Sprintf("%v", play.Spec.ProjectID),
	})
	envVars = append(envVars, corev1.EnvVar{
		Name:  config.GetEnvPrefix() + "PIPELINE_ID",
		Value: fmt.Sprintf("%v", play.Spec.PipelineID),
	})
	envVars = append(envVars, corev1.EnvVar{
		Name:  config.GetEnvPrefix() + "REPOSITORY_URL",
		Value: play.Spec.RepoURL,
	})
	envVars = append(envVars, corev1.EnvVar{
		Name:  config.GetEnvPrefix("commit") + "SHA",
		Value: play.Spec.Commit.SHA,
	})
	envVars = append(envVars, corev1.EnvVar{
		Name:  config.GetEnvPrefix("commit") + "BEFORE_SHA",
		Value: play.Spec.Commit.BeforeSHA,
	})
	if len(play.Spec.Commit.SHA) >= 8 {
		envVars = append(envVars, corev1.EnvVar{
			Name:  config.GetEnvPrefix("commit") + "SHORT_SHA",
			Value: play.Spec.Commit.SHA[:8],
		})
	}
	envVars = append(envVars, corev1.EnvVar{
		Name:  config.GetEnvPrefix("commit") + "REF_NAME",
		Value: play.Spec.Commit.Ref,
	})
	envVars = append(envVars, corev1.EnvVar{
		Name:  config.GetEnvPrefix("expose") + "DOMAIN",
		Value: play.Spec.Domain,
	})
	envVars = append(envVars, corev1.EnvVar{
		Name:  config.GetEnvPrefix() + "EXPOSE",
		Value: strconv.FormatBool(play.Spec.Expose),
	})
	envVars = append(envVars, corev1.EnvVar{
		Name:  config.GetEnvPrefix() + "EXTERNAL",
		Value: strconv.FormatBool(play.Spec.External),
	})
	envVars = append(envVars, corev1.EnvVar{
		Name:  config.GetEnvPrefix() + "DOCKER_URL",
		Value: play.Spec.DockerURL,
	})

	return
}

// BuildTaskPredefinedEnv build and return the environment variables from tasks
func BuildTaskPredefinedEnv(tasks []map[ci.TaskType]ci.Task) (envVars []corev1.EnvVar) {

	for pos := range tasks {
		for name, task := range tasks[pos] {
			prefix := config.GetEnvPrefix(string(name))
			envVars = append(envVars, corev1.EnvVar{
				Name:  prefix + "CONTEXT",
				Value: task.Docker.Context,
			})
			envVars = append(envVars, corev1.EnvVar{
				Name:  prefix + "DOCKERFILE",
				Value: task.Docker.Filepath,
			})
			envVars = append(envVars, corev1.EnvVar{
				Name:  prefix + "IMAGE",
				Value: task.Image,
			})
			envVars = append(envVars, corev1.EnvVar{
				Name:  prefix + "NAMESPACE",
				Value: task.Namespace,
			})
		}
	}

	return
}

// BuildConfigPredefinedEnv return then envVar from config values
func BuildConfigPredefinedEnv() (envVars []corev1.EnvVar) {
	cfg := config.GetConfig()
	prefix := config.GetEnvPrefix("config")
	envVars = append(envVars, corev1.EnvVar{
		Name:  prefix + "DEFAULT_DOMAIN",
		Value: cfg.DefaultDomain,
	})
	envVars = append(envVars, corev1.EnvVar{
		Name:  prefix + "CLUSTER_ROLE",
		Value: cfg.ClusterRole,
	})
	envVars = append(envVars, corev1.EnvVar{
		Name:  prefix + "INGRESS_CLASS",
		Value: cfg.Ingress.Class,
	})
	envVars = append(envVars, corev1.EnvVar{
		Name:  prefix + "NAMESPACE",
		Value: cfg.Namespace,
	})

	return
}
