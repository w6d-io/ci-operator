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

package secrets

import (
	"context"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/k8s/sa"
	"github.com/w6d-io/ci-operator/internal/util"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	DockerPrefixSecret = "reg-cred"
)

// DockerCredCreate creates the docker config json secret and add it into the service account
func (s *Secret) DockerCredCreate(ctx context.Context, r client.Client, log logr.Logger) error {
	log = log.WithName("Create").WithValues("action", DockerPrefixSecret)
	log.V(1).Info("creating")

	namespacedName := util.GetCINamespacedName(DockerPrefixSecret, s.Play)
	resource := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        namespacedName.Name,
			Namespace:   namespacedName.Namespace,
			Annotations: make(map[string]string),
			Labels:      util.GetCILabels(s.Play),
		},
		StringData: map[string]string{
			corev1.DockerConfigJsonKey: s.GetSecret(ci.DockerConfig, log),
		},
		Type: corev1.SecretTypeDockerConfigJson,
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
