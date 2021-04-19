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
Created on 16/12/2020
*/

package pipelinerun

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/k8s/sa"
	"github.com/w6d-io/ci-operator/internal/k8s/secrets"
	"github.com/w6d-io/ci-operator/internal/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (p *PipelineRun) Parse(log logr.Logger) error {

	for pos, m := range p.Play.Spec.Tasks {
		for name := range m {
			switch name {
			case ci.UnitTests:
				_ = p.SetUnitTest(pos, log)
			case ci.Build:
				if err := p.SetBuild(pos, log); err != nil {
					return err
				}
			case ci.Sonar:
				_ = p.SetSonar(pos, log)
			case ci.Deploy:
				_ = p.SetDeploy(pos, log)
			case ci.IntegrationTests:
				_ = p.SetIntTest(pos, log)
			case ci.Clean:
				_ = p.SetClean(pos, log)
			case ci.E2ETests:
				_ = p.SetE2ETest(pos, log)
			}
		}
	}
	p.PodTemplate = config.PodTemplate().DeepCopy()
	if config.GetMinio().Host != "" {
		p.PodTemplate.Volumes = append(p.PodTemplate.Volumes, corev1.Volume{
			Name: secrets.MinIOPrefixSecret,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: util.GetCINamespacedName(secrets.MinIOPrefixSecret, p.Play).Name,
				},
			},
		})
	}
	var okVault, okSecret bool
	if p.Play.Spec.Vault != nil {
		_, okVault = p.Play.Spec.Vault.Secrets[secrets.KubeConfigKey]
	}
	_, okSecret = p.Play.Spec.Secret[secrets.KubeConfigKey]
	if okVault || okSecret {
		p.PodTemplate.Volumes = append(p.PodTemplate.Volumes, corev1.Volume{
			Name: secrets.KubeConfigPrefix,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: util.GetCINamespacedName(secrets.KubeConfigPrefix, p.Play).Name,
				},
			},
		})
	}
	return nil
}

func (p *PipelineRun) Create(ctx context.Context, r client.Client, log logr.Logger) error {
	log = log.WithName("Create").WithValues("action", "pipeline-run")
	log.V(1).Info("creating")
	namespacedName := util.GetCINamespacedName(Prefix, p.Play)
	resource := &tkn.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:        namespacedName.Name,
			Namespace:   namespacedName.Namespace,
			Annotations: make(map[string]string),
			Labels:      util.GetCILabels(p.Play),
		},
		Spec: tkn.PipelineRunSpec{
			PipelineRef:        &tkn.PipelineRef{Name: util.GetCINamespacedName("pipeline", p.Play).Name},
			ServiceAccountName: util.GetCINamespacedName(sa.Prefix, p.Play).Name,
			Resources:          p.getPipelineResourceBinding(p.Play),
			Workspaces:         p.getWorkspaceBinding(),
			Params:             p.Params,
			PodTemplate:        p.PodTemplate,
		},
	}
	resource.Annotations[config.ScheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
	if err := controllerutil.SetControllerReference(p.Play, resource, p.Scheme); err != nil {
		return err
	}
	log.V(1).Info(resource.Kind, "content", fmt.Sprintf("%v",
		util.GetObjectContain(resource)))
	if err := r.Create(ctx, resource); err != nil {
		return err
	}
	// All went well
	return nil
}

func (p *PipelineRun) getWorkspaceBinding() []tkn.WorkspaceBinding {
	return []tkn.WorkspaceBinding{
		config.Volume(),
	}
}

func (p *PipelineRun) getPipelineResourceBinding(play *ci.Play) []tkn.PipelineResourceBinding {
	res := []tkn.PipelineResourceBinding{
		{
			Name: ci.ResourceGit,
			ResourceRef: &tkn.PipelineResourceRef{
				Name: util.GetCINamespacedName("pr-git", play).Name,
			},
		},
	}
	if util.IsBuildStage(play) {
		res = append(res, tkn.PipelineResourceBinding{
			Name: ci.ResourceImage,
			ResourceRef: &tkn.PipelineResourceRef{
				Name: util.GetCINamespacedName("pr-image", play).Name,
			},
		})
	}
	return res
}
