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
Created on 28/12/2020
*/

package sa

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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// CI structure for build ServiceAccount k8s resource
type CI struct {
	internal.WorkFlowStruct
}

const (
	// Prefix contains the prefix for CI resources
	Prefix = "sa"
)

func (s *CI) Create(ctx context.Context, r client.Client, log logr.Logger) error {
	log = log.WithName("Create").WithValues("action", Prefix)
	log.V(1).Info("creating")

	namespacedName := util.GetCINamespacedName(Prefix, s.Play)
	log.V(1).WithValues("namespaced", namespacedName).Info("debug")
	resource := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:        namespacedName.Name,
			Namespace:   namespacedName.Namespace,
			Annotations: make(map[string]string),
			Labels:      util.GetCILabels(s.Play),
		},
	}
	resource.Annotations[config.ScheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
	if err := controllerutil.SetControllerReference(s.Play, resource, s.Scheme); err != nil {
		return err
	}
	log.V(1).Info(resource.Kind, "content", fmt.Sprintf("%v",
		util.GetObjectContain(resource)))

	if err := r.Create(ctx, resource); err != nil {
		log.Error(err, "create")
		return err
	}
	return nil
}
