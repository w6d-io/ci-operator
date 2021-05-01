/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 21/04/2021
*/

package secrets_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal"
	"github.com/w6d-io/ci-operator/internal/k8s/secrets"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

var _ = Describe("", func() {
	Context("", func() {
		It("validates all steps until succeeded", func() {
			var err error
			s := &secrets.Secret{
				WorkFlowStruct: internal.WorkFlowStruct{
					Play: &ci.Play{
						Spec: ci.PlaySpec{
							ProjectID:  13,
							PipelineID: 1,
							Secret: map[ci.SecretKind]string{
								ci.DockerConfig: "{}",
							},
						},
					},
				},
			}
			By("fail controller reference")
			s.Scheme = runtime.NewScheme()
			err = s.MinIOCreate(ctx, k8sClient, ctrl.Log)
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(ContainSubstring("no kind is registered for the type"))

			By("fail to create")
			s.Scheme = scheme
			err = s.MinIOCreate(ctx, k8sClient, ctrl.Log)
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(Equal(`namespaces "p6e-cx-13" not found`))

			By("create namespace")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "p6e-cx-13",
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(Succeed())

			By("create failed on sa update")
			s.Play.Name = "play-test-13-1"
			s.Play.UID = "77557df7-3162-46cf-9b3c-d3b9b70a42b8"
			err = s.MinIOCreate(ctx, k8sClient, ctrl.Log)
			Expect(err).To(Succeed())

		})
	})
})
