/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 12/01/2021
*/

package rbac

import (
	"context"
	"fmt"
	"time"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/k8s/serviceaccount"
	"github.com/w6d-io/ci-operator/internal/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CI structure for build Rbac k8s resource
type Deploy struct {
	internal.WorkFlowStruct
}

func (in *Deploy) Create(ctx context.Context, r client.Client, logger logr.Logger) error {
	log := logger.WithName("Create").WithValues("action", Prefix)
	log.V(1).Info("creating")

	namespacedName := util.GetCINamespacedName(Prefix, in.Play)
	deployNamespacedNamed := util.GetDeployNamespacedName(config.GetDeployPrefix(), in.Play)
	log.V(1).WithValues("namespaced", namespacedName).Info("debug")
	resource := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:        namespacedName.Name,
			Namespace:   deployNamespacedNamed.Namespace,
			Annotations: make(map[string]string),
			Labels:      util.GetCILabels(in.Play),
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "ClusterRole",
			Name:     config.GetClusterRole(),
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      util.GetCINamespacedName(serviceaccount.Prefix, in.Play).Name,
				Namespace: namespacedName.Namespace,
			},
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      util.GetCINamespacedName(serviceaccount.Prefix, in.Play).Name,
				Namespace: deployNamespacedNamed.Namespace,
			},
		},
	}
	resource.Annotations[config.ScheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
	//if err := controllerutil.SetControllerReference(in.Play, resource, in.Scheme); err != nil {
	//	return err
	//}
	log.V(1).Info(fmt.Sprintf("rolbinding contains\n%v",
		util.GetObjectContain(resource)))
	if err := r.Create(ctx, resource); err != nil {
		return err
	}
	return nil
}
