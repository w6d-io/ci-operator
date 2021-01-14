/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 28/12/2020
*/

package serviceaccount

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Deploy structure for build ServiceAccount k8s resource
type Deploy struct {
	internal.WorkFlowStruct
}

func (s *Deploy) Create(ctx context.Context, r client.Client, log logr.Logger) error {
	log = log.WithName("Create").WithValues("action", Prefix)
	log.V(1).Info("creating")

	namespacedName := util.GetDeployNamespacedName(Prefix, s.Play)
	log.V(1).WithValues("namespaced", namespacedName).Info("debug")
	resource := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:        namespacedName.Name,
			Namespace:   util.GetDeployNamespacedName(config.GetDeployPrefix(), s.Play).Namespace,
			Annotations: make(map[string]string),
			Labels:      util.GetCILabels(s.Play),
		},
	}
	resource.Annotations[config.ScheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
	log.V(1).Info(fmt.Sprintf("Secret contains\n%v",
		util.GetObjectContain(resource)))
	if err := r.Create(ctx, resource); err != nil {
		return err
	}

	return nil
}
