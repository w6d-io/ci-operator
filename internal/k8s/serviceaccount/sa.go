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

package serviceaccount

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Update(ctx context.Context, name string, nn types.NamespacedName, r client.Client) error {
	// Add secret in CI
	saResource := new(corev1.ServiceAccount)
	if err := r.Get(ctx, nn, saResource); err != nil {
		return err
	}
	if !isContainSecret(name, saResource.Secrets) {
		saResource.Secrets = append(saResource.Secrets,
			corev1.ObjectReference{Name: name})
		if err := r.Update(ctx, saResource); err != nil {
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
