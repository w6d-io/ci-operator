/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 29/04/2021
*/

package sa_test

import (
	"github.com/w6d-io/ci-operator/internal/k8s/sa"
	"k8s.io/apimachinery/pkg/types"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service account", func() {
	Context("", func() {
		BeforeEach(func() {
		})
		AfterEach(func() {
		})
		It("update service account", func() {
			var err error
			err = sa.Update(ctx, "test-secret", types.NamespacedName{
				Namespace: "p6e-cx-8",
				Name:      "sa-8-1",
			}, k8sClient)
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(Equal(`serviceaccounts "sa-8-1" not found`))

			By("create namespace")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "p6e-cx-8",
				},
			}
			err = k8sClient.Create(ctx, ns)
			Expect(err).To(Succeed())

			By("create service account")
			s := &corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sa-8-1",
					Namespace: "p6e-cx-8",
				},
			}
			err = k8sClient.Create(ctx, s)
			Expect(err).To(Succeed())

			By("update service account with namespace and service account exists")
			err = sa.Update(ctx, "test-secret", types.NamespacedName{
				Namespace: "p6e-cx-8",
				Name:      "sa-8-1",
			}, k8sClient)
			Expect(err).To(Succeed())

			By("update service account with namespace and service account exists and secret already in")
			err = sa.Update(ctx, "test-secret", types.NamespacedName{
				Namespace: "p6e-cx-8",
				Name:      "sa-8-1",
			}, k8sClient)
			Expect(err).To(Succeed())

		})
	})
})
