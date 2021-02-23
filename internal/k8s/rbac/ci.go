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
Created on 29/12/2020
*/

package rbac

import (
	"context"
	"fmt"
	"github.com/w6d-io/ci-operator/internal/k8s/sa"
	"time"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// CI structure for build Rbac k8s resource
type CI struct {
	internal.WorkFlowStruct
}

const (
	// Prefix contains the prefix for Rbac resources
	Prefix = "role-binding"
)

func (in *CI) Create(ctx context.Context, r client.Client, logger logr.Logger) error {
	log := logger.WithName("Create").WithValues("action", Prefix)
	log.V(1).Info("creating")

	namespacedName := util.GetCINamespacedName(Prefix, in.Play)
	log.V(1).WithValues("namespaced", namespacedName).Info("debug")
	resource := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:        namespacedName.Name,
			Namespace:   namespacedName.Namespace,
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
				Name:      util.GetCINamespacedName(sa.Prefix, in.Play).Name,
				Namespace: namespacedName.Namespace,
			},
		},
	}
	resource.Annotations[config.ScheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
	if err := controllerutil.SetControllerReference(in.Play, resource, in.Scheme); err != nil {
		return err
	}
	log.V(1).Info(fmt.Sprintf("rolbinding contains\n%v",
		util.GetObjectContain(resource)))
	if err := r.Create(ctx, resource); err != nil {
		return err
	}

	return nil
}
