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
Created on 31/03/2021
*/
package configmap_test

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/k8s/configmap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Configmap", func() {
	Context("call get content", func() {
		BeforeEach(func() {
		})
		AfterEach(func() {
		})
		It("get empty string with nil parameter", func() {
			str := configmap.GetContentFromKeySelector(nil, nil, nil)
			Expect(str).To(Equal(""))
		})
		It("", func() {
			var err error
			config.SetNamespace("default")
			ctx := context.WithValue(context.Background(), "correlation_id", "unit-test")
			cm := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "unit-test",
					Namespace: config.GetNamespace(),
				},
				Data: map[string]string{
					"values.yaml": "Test",
				},
			}
			err = k8sClient.Create(ctx, cm)
			Expect(err).To(Succeed())
			c := &corev1.ConfigMapKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "unit-test",
				},
				Key: "nonexistent",
			}
			str := configmap.GetContentFromKeySelector(context.Background(), k8sClient, c)
			Expect(str).To(Equal(""))
			err = k8sClient.Delete(ctx, cm)
			Expect(err).To(Succeed())
		})
	})
})
