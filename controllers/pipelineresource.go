/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 22/11/2020
*/

package controllers

import (
	"context"
	"errors"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"net/url"
	"time"

	resourcev1alpha1 "github.com/tektoncd/pipeline/pkg/apis/resource/v1alpha1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GitPR git PipelineResource type for CI
type GitPR struct {
	OwnerReferences []metav1.OwnerReference
	NamespacedName  types.NamespacedName
	Labels          map[string]string
	Revision        string
	URL             *url.URL
}

// Create implements the CIInterface method
func (g GitPR) Create(ctx context.Context, r client.Client, log logr.Logger) error {
	log = log.WithName("Create").WithValues("pipelineResource", "git")
	log.V(1).Info("create git pipelineResource")
	gitResource := &resourcev1alpha1.PipelineResource{
		ObjectMeta: metav1.ObjectMeta{
			Name:            g.NamespacedName.Name,
			Namespace:       g.NamespacedName.Namespace,
			Annotations:     make(map[string]string),
			Labels:          g.Labels,
			OwnerReferences: g.OwnerReferences,
		},
		Spec: resourcev1alpha1.PipelineResourceSpec{
			Type: resourcev1alpha1.PipelineResourceTypeGit,
			Params: []resourcev1alpha1.ResourceParam{
				{
					Name:  "revision",
					Value: g.Revision,
				},
				{
					Name:  "URL",
					Value: g.URL.String(),
				},
			},
		},
	}

	// set the current time in the new pipeline resource git type resource in annotation
	gitResource.Annotations[scheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
	if err := r.Create(ctx, gitResource); err != nil {
		return err
	}
	return nil
}

// SetGitCreate adds git pipeline resource
func (wf *WFType) SetGitCreate(p ci.Play, r *PlayReconciler) error {
	log := r.Log.WithName("SetGitCreate").WithValues("cx-namespace", InNamespace(p))
	log.V(1).Info("Check repository URL")
	URL, err := url.Parse(p.Spec.RepoURL)
	if err != nil {
		log.Error(err, "URL parse error")
		return err
	}
	git := &GitPR{
		// TODO put the prefix in config
		NamespacedName:  CxCINamespacedName("pr-git", p),
		URL:             URL,
		Labels:          CxCILabels(p),
		Revision:        p.Spec.Commit.Ref,
		OwnerReferences: []metav1.OwnerReference{CIOwnerReference(p)},
	}
	if err := wf.Add(git.Create); err != nil {
		return err
	}
	// TODO add git token secret creation
	log.V(1).Info("GitCreate added")
	return nil
}

// ImagePR image PipelineResource type for CI
type ImagePR struct {
	OwnerReferences []metav1.OwnerReference
	NamespacedName  types.NamespacedName
	Labels          map[string]string
	URL             *url.URL
}

// Create implements the CIInterface method
func (i *ImagePR) Create(ctx context.Context, r client.Client, log logr.Logger) error {
	log = log.WithName("Create").WithValues("pipelineResource", "image")
	log.V(1).Info("create image pipelineResource")
	imageResource := &resourcev1alpha1.PipelineResource{
		ObjectMeta: metav1.ObjectMeta{
			Name:        i.NamespacedName.Name,
			Namespace:   i.NamespacedName.Namespace,
			Annotations: make(map[string]string),
			Labels:      i.Labels,
		},
		Spec: resourcev1alpha1.PipelineResourceSpec{
			Type: resourcev1alpha1.PipelineResourceTypeImage,
			Params: []resourcev1alpha1.ResourceParam{
				{
					Name:  "URL",
					Value: i.URL.String(),
				},
			},
		},
	}

	// set the current time in the new pipeline resource image type resource in annotation
	imageResource.Annotations[scheduledTimeAnnotation] = time.Now().Format(time.RFC3339)

	if err := r.Create(ctx, imageResource); err != nil {
		if apierrors.IsAlreadyExists(err) {
			log.Error(err, "creating failed")
			return nil
		}
		return err
	}
	return nil
}

// SetImageCreate adds image pipeline resource
func (wf *WFType) SetImageCreate(p ci.Play, r *PlayReconciler) error {
	if IsBuildStage(p) {
		log := r.Log.WithName("SetImageCreate").WithValues("cx-namespace", InNamespace(p))
		URL, err := CxDockerImageName(p)
		if err != nil {
			log.Error(err, "get repository address failed")
			return errors.New("get repository address failed")
		}
		image := ImagePR{
			// TODO set the prefix in configuration
			NamespacedName: CxCINamespacedName("pr-image", p),
			// TODO allow
			URL:    URL,
			Labels: CxCILabels(p),
		}
		if err := wf.Add(image.Create); err != nil {
			return err
		}
		// TODO add registry credential secret creation
		log.V(1).Info("ImageCreate added")

	}
	return nil
}
