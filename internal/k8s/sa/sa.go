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

package sa

import (
	"context"
	"fmt"
	"github.com/w6d-io/ci-operator/internal/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Update adds the secret reference in the service account
func Update(ctx context.Context, name string, nn types.NamespacedName, r client.Client) error {
	correlationID := ctx.Value("correlation_id")
	log := ctrl.Log.WithValues("correlation_id", correlationID, "action", "service-account").
		WithName("controllers").
		WithName("Play").
		WithName("CreateCI").
		WithName("Create").
		WithName("Update").
		WithValues("name", name, "object", nn.String())
	log.V(1).Info("update")
	resource := new(corev1.ServiceAccount)
	if err := r.Get(ctx, nn, resource); err != nil {
		log.Error(err, "fail to get service account")
		return err
	}
	if !isContainSecret(name, resource.Secrets) {
		log.V(1).Info("append secret in sa")
		err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			log.V(1).Info("get sa")
			if err := r.Get(ctx, nn, resource); err != nil {
				log.Error(err, "fail to get service account (retry)")
				return err
			}
			resource.Secrets = append(resource.Secrets,
				corev1.ObjectReference{Name: name})
			log.V(1).Info("update sa", "content", fmt.Sprintf("%v",
				util.GetObjectContain(resource)))
			if err := r.Update(ctx, resource); err != nil {
				log.Error(err, "fail to update service account (retry)")
				return err
			}
			return nil
		})
		if err != nil {
			log.Error(err, "fail to update service account")
			return err
		}
	}
	return nil
}

func isContainSecret(name string, secrets []corev1.ObjectReference) bool {
	for _, secret := range secrets {
		if secret.Name == name {
			return true
		}
	}
	return false
}
