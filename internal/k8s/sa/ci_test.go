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
	"github.com/w6d-io/ci-operator/internal"
	"github.com/w6d-io/ci-operator/internal/k8s/sa"
	"k8s.io/apimachinery/pkg/runtime"

	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service Account", func() {
	Context("Create", func() {
		BeforeEach(func() {
		})
		AfterEach(func() {
		})
		It("failed because no kind", func() {
			var err error

			s := &sa.CI{
				WorkFlowStruct: internal.WorkFlowStruct{
					Play: &ci.Play{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "play-unit-test",
							Namespace: "default",
						},
						Spec: ci.PlaySpec{
							ProjectID:  6,
							PipelineID: 1,
						},
					},
					Scheme: runtime.NewScheme(),
				},
			}
			err = s.Create(ctx, k8sClient, ctrl.Log)
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(ContainSubstring("no kind is registered for the type v1alpha1.Play"))
		})
		It("failed because namespace does not exist", func() {
			var err error

			s := &sa.CI{
				WorkFlowStruct: internal.WorkFlowStruct{
					Play: &ci.Play{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "play-unit-test",
							Namespace: "p6e-cx-6",
						},
						Spec: ci.PlaySpec{
							ProjectID:  6,
							PipelineID: 1,
						},
					},
					Scheme: scheme,
				},
			}
			err = s.Create(ctx, k8sClient, ctrl.Log)
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(Equal(`namespaces "p6e-cx-6" not found`))
		})
		It("success", func() {
			var err error

			By("create the namespace")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "p6e-cx-6",
				},
			}
			err = k8sClient.Create(ctx, ns)
			Expect(err).To(Succeed())
			s := &sa.CI{
				WorkFlowStruct: internal.WorkFlowStruct{
					Play: &ci.Play{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "play-unit-test",
							Namespace: "p6e-cx-6",
							UID:       "ca0673b2-d121-4d2a-8d49-244ac0a82d72",
						},
						Spec: ci.PlaySpec{
							ProjectID:  6,
							PipelineID: 1,
						},
					},
					Scheme: scheme,
				},
			}
			err = s.Create(ctx, k8sClient, ctrl.Log)
			Expect(err).To(Succeed())
		})
	})
})
