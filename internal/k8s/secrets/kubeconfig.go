/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 16/04/2021
*/

package secrets

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/k8s/sa"
	"github.com/w6d-io/ci-operator/internal/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

const (
	KubeConfigKey    = "kubeconfig"
	KubeConfigPrefix = "kubeconfig"
)

func (s *Secret) KubeConfigCreate(ctx context.Context, r client.Client, logger logr.Logger) error {
	log := logger.WithName("Create").WithValues("action", KubeConfigPrefix)
	log.V(1).Info("creating")

	namespacedName := util.GetCINamespacedName(KubeConfigPrefix, s.Play)
	resource := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        namespacedName.Name,
			Namespace:   namespacedName.Namespace,
			Labels:      util.GetCILabels(s.Play),
			Annotations: make(map[string]string),
		},
		StringData: map[string]string{
			"config": s.GetSecret(KubeConfigKey, log),
		},
		Type: corev1.SecretTypeOpaque,
	}

	resource.Annotations[config.ScheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
	if err := controllerutil.SetControllerReference(s.Play, resource, s.Scheme); err != nil {
		return err
	}
	if err := r.Create(ctx, resource); err != nil {
		log.Error(err, "create")
		return err
	}
	if err := sa.Update(ctx, resource.Name,
		util.GetCINamespacedName(sa.Prefix, s.Play), r); err != nil {
		return err
	}
	return nil
}
