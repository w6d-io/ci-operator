/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 16/12/2020
*/

package pipelinerun

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"time"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (p *PipelineRun) Parse(play ci.Play, log logr.Logger) func(context.Context, client.Client, logr.Logger) error {

	for pos, m := range play.Spec.Tasks {
		for name := range m {
			switch name {
			case ci.Build:
				if err := p.SetBuild(pos, play, log); err != nil {
					return nil
				}
				//		case ci.Sonar:
				//			if err := SetPipelineRunSonar(play, r); err != nil {
				//				return err
				//			}
				//		case ci.UnitTests:
				//			if err := SetPipelineRunUnitTest(play, r); err != nil {
				//				return err
				//			}
				//		case ci.IntegrationTests:
				//			if err := SetPipelineRunIntTest(play, r); err != nil {
				//				return err
				//			}
				//		case ci.Deploy:
				//			if err := SetPipelineRunDeploy(play, r); err != nil {
				//				return err
				//			}
				//		case ci.Clean:
				//			if err := SetPipelineRunClean(play, r); err != nil {
				//				return err
				//			}
			case ci.Clean:
				if err := p.SetClean(pos, play, log); err != nil {
					return nil
				}
			case ci.Deploy:
				if err := p.SetDeploy(pos, play, log); err != nil {
					return nil
				}
			case ci.IntegrationTests:
				if err := p.SetIntTest(pos, play, log); err != nil {
					return nil
				}
			case ci.Sonar:
				if err := p.SetSonar(pos, play, log); err != nil {
					return nil
				}
			case ci.UnitTests:
				if err := p.SetUnitTest(pos, play, log); err != nil {
					return nil
				}
			}
		}
	}
	return p.Create
}

func (p *PipelineRun) Create(ctx context.Context, r client.Client, log logr.Logger) error {
	log = log.WithName("Create").WithValues("action", "pipeline-run")
	log.V(1).Info("creating")
	namespacedName := p.GetNamespacedName("pipeline-run", p.Play)
	pipelineRunResource := &tkn.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:            namespacedName.Name,
			Namespace:       namespacedName.Namespace,
			Annotations:     make(map[string]string),
			Labels:          p.GetLabels(p.Play),
			OwnerReferences: []metav1.OwnerReference{p.GetOwnerReferences(p.Play)},
		},
		Spec: tkn.PipelineRunSpec{
			PipelineRef:        &tkn.PipelineRef{Name: p.GetNamespacedName("pipeline", p.Play).Name},
			ServiceAccountName: p.ServiceAccount,
			Resources:          p.getPipelineResourceBinding(p.Play),
			Workspaces:         p.getWorkspaceBinding(),
			Params:             p.Params,
		},
	}
	log.V(1).Info(fmt.Sprintf("pipelineRun contains\n%v",
		p.GetObjectContain(pipelineRunResource)))
	pipelineRunResource.Annotations[p.scheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
	if err := r.Create(ctx, pipelineRunResource); err != nil {
		return err
	}
	// All went well
	return nil
}

func (p *PipelineRun) getWorkspaceBinding() []tkn.WorkspaceBinding {
	res := corev1.ResourceList{}
	res[corev1.ResourceStorage] = resource.MustParse("1Gi")
	mode := corev1.ReadWriteOnce
	if p.Volume.Mode != "" {
		mode = p.Volume.Mode
	}
	if p.Volume.Size != "" {
		res[corev1.ResourceStorage] = resource.MustParse(p.Volume.Size)
	}
	var wks []tkn.WorkspaceBinding
	for _, volume := range p.WorkspacesDeclaration {
		wks = append(wks, tkn.WorkspaceBinding{
			Name: volume.Name,
			VolumeClaimTemplate: &corev1.PersistentVolumeClaim{
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{mode},
					Resources: corev1.ResourceRequirements{
						Requests: res,
					},
				},
			},
		})
	}
	return wks
}

func (p *PipelineRun) getPipelineResourceBinding(play ci.Play) []tkn.PipelineResourceBinding {
	res := []tkn.PipelineResourceBinding{
		{
			Name: ci.ResourceGit,
			ResourceRef: &tkn.PipelineResourceRef{
				Name: p.GetNamespacedName("pr-git", play).Name,
			},
		},
	}
	if p.IsBuildStage(play) {
		res = append(res, tkn.PipelineResourceBinding{
			Name: ci.ResourceImage,
			ResourceRef: &tkn.PipelineResourceRef{
				Name: p.GetNamespacedName("pr-image", play).Name,
			},
		})
	}
	return res
}
