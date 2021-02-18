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
Created on 12/01/2021
*/

package rbac

import (
	"context"
	"fmt"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"time"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/k8s/sa"
	"github.com/w6d-io/ci-operator/internal/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CI structure for build Rbac k8s resource
type Deploy struct {
	internal.WorkFlowStruct
}

func (in *Deploy) Create(ctx context.Context, r client.Client, logger logr.Logger) error {
	log := logger.WithName("Deploy").WithName("Create").WithValues("action", Prefix)
	log.V(1).Info("creating")

	namespacedNamed := util.GetCINamespacedName2(Prefix, in.Play)
	deployNamespacedNamed := util.GetDeployNamespacedName(config.GetDeployPrefix(), in.Play)
	log.V(1).WithValues("namespaced", namespacedNamed).Info("debug")
	resource := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:        namespacedNamed.Name,
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
				Name:      util.GetCINamespacedName2(sa.Prefix, in.Play).Name,
				Namespace: namespacedNamed.Namespace,
			},
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      util.GetDeployNamespacedName(sa.Prefix, in.Play).Name,
				Namespace: deployNamespacedNamed.Namespace,
			},
		},
	}
	resource.Annotations[config.ScheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
	// TODO find a way to link this resource with oper
	log.V(1).Info(fmt.Sprintf("rolebinding contains\n%v",
		util.GetObjectContain(resource)))
	old := &rbacv1.RoleBinding{}
	err := r.Get(ctx, types.NamespacedName{}, old)
	if apierrors.IsNotFound(err) {
		if err := r.Create(ctx, resource); err != nil {
			log.Error(err, "create")
			return err
		}
		return nil
	}
	if err := r.Update(ctx, resource); err != nil {
		log.Error(err, "update")
		return err
	}
	return nil
}
