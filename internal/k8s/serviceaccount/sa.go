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
