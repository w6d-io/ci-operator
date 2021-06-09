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
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
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

// Deploy structure for build Rbac k8s resource
type Deploy struct {
	internal.WorkFlowStruct
}

func (in *Deploy) Create(ctx context.Context, r client.Client, logger logr.Logger) error {
	log := logger.WithName("Deploy").WithName("Create").WithValues("action", Prefix)
	log.V(1).Info("creating")

	if in.Play.IsInternal() {
		log.V(1).Info("skip")
		return nil
	}
	namespacedNamed := util.GetCINamespacedName2(Prefix, in.Play)
	deployNamespacedNamed := util.GetDeployNamespacedName(config.GetDeployPrefix(), in.Play)

	// TODO find a way to link this resource with play
	resource := &rbacv1.RoleBinding{}
	nn := types.NamespacedName{Name: namespacedNamed.Name, Namespace: deployNamespacedNamed.Namespace}
	err := r.Get(ctx, nn, resource)
	if apierrors.IsNotFound(err) {
		log.V(1).Info("Create")
		resource = GetRoleBinding(in.Play)
		resource.Annotations[config.ScheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
		if err := r.Create(ctx, resource); err != nil {
			log.Error(err, "create")
			return err
		}
		return nil
	}
	log.V(1).Info("Update")
	if isSubjectExist(GetSubject(in.Play), resource.Subjects) {
		return nil
	}
	err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		if err := r.Get(ctx, nn, resource); err != nil {
			return err
		}
		resource.Subjects = append(resource.Subjects, GetSubject(in.Play))
		if err := r.Update(ctx, resource); err != nil {
			log.Error(err, "update")
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// GetRoleBinding return the a new role binding resource
func GetRoleBinding(play *ci.Play) *rbacv1.RoleBinding {
	namespacedNamed := util.GetCINamespacedName2(Prefix, play)
	deployNamespacedNamed := util.GetDeployNamespacedName(config.GetDeployPrefix(), play)
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:        namespacedNamed.Name,
			Namespace:   deployNamespacedNamed.Namespace,
			Annotations: make(map[string]string),
			Labels:      util.GetCILabels(play),
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "ClusterRole",
			Name:     config.GetClusterRole(),
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      util.GetCINamespacedName(sa.Prefix, play).Name,
				Namespace: namespacedNamed.Namespace,
			},
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      util.GetDeployNamespacedName(sa.Prefix, play).Name,
				Namespace: deployNamespacedNamed.Namespace,
			},
		},
	}
}

// GetSubject return a rbac role binding subject element
func GetSubject(play *ci.Play) rbacv1.Subject {
	return rbacv1.Subject{
		Kind:      rbacv1.ServiceAccountKind,
		Name:      util.GetCINamespacedName(sa.Prefix, play).Name,
		Namespace: util.GetCINamespacedName2(Prefix, play).Namespace,
	}
}

// isSubjectExists check the subject presence
func isSubjectExist(needle rbacv1.Subject, haystack []rbacv1.Subject) bool {
	for _, value := range haystack {
		if needle.Name == value.Name && needle.Namespace == value.Namespace {
			return true
		}
	}
	return false
}
