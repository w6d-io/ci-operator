/*
Copyright 2020 WILDCARD SA.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
Created on 30/03/2021
*/
package configmap

import (
	"context"
	"github.com/w6d-io/ci-operator/internal/config"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var logger = ctrl.Log

func GetContentFromKeySelector(ctx context.Context, r client.Client, c *corev1.ConfigMapKeySelector) string {
	if r == nil || c == nil {
		logger.V(1).Info("k8s client or configmap key is nil")
		return ""
	}
	cm := &corev1.ConfigMap{}
	o := client.ObjectKey{Name: c.Name, Namespace: config.GetNamespace()}
	err := r.Get(ctx, o, cm)
	if err != nil {
		logger.Error(err, "get configmap", "name", c.Name)
		return ""
	}
	content, ok := cm.Data[c.Key]
	if !ok {
		logger.Error(nil, "no such element in configmap", "name", c.Name, "key", c.Key)
		return ""
	}
	return content
}
