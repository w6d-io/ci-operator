/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 19/04/2021
*/

package v1alpha1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/w6d-io/ci-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Webhook", func() {
	Context("Default", func() {
		BeforeEach(func() {
		})
		AfterEach(func() {
		})

	})
	Context("Create", func() {
		It("success", func() {
			var err error
			p := &v1alpha1.Play{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: v1alpha1.PlaySpec{
					Scope: v1alpha1.Scope{
						Name: "",
					},
					Environment: "develop",
					ProjectID:   1,
					PipelineID:  1,
					RepoURL:     "https://github.com",
					Commit: v1alpha1.Commit{
						SHA: "sha",
						Ref: "main",
					},
					Tasks: []map[v1alpha1.TaskType]v1alpha1.Task{
						{
							v1alpha1.Deploy: v1alpha1.Task{},
						},
					},
				},
			}
			err = p.ValidateCreate()
			Expect(err).To(Succeed())
		})
		It("fails on wrong test", func() {
			var err error
			p := &v1alpha1.Play{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: v1alpha1.PlaySpec{
					Scope: v1alpha1.Scope{
						Name: "",
					},
					Tasks: []map[v1alpha1.TaskType]v1alpha1.Task{
						{
							"test": v1alpha1.Task{},
						},
						{
							v1alpha1.Deploy: v1alpha1.Task{},
						},
					},
				},
			}
			err = p.ValidateCreate()
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(ContainSubstring("Invalid value: test"))
		})
		It("fails on domain and repo_url", func() {
			var err error
			p := &v1alpha1.Play{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: v1alpha1.PlaySpec{
					Scope: v1alpha1.Scope{
						Name: "",
					},
					RepoURL: "http://{}",
					Domain:  "ee_@",
					Tasks: []map[v1alpha1.TaskType]v1alpha1.Task{
						{
							"test": v1alpha1.Task{},
						},
						{
							v1alpha1.Deploy: v1alpha1.Task{},
						},
					},
					Vault: &v1alpha1.Vault{
						Secrets: map[v1alpha1.SecretKind]v1alpha1.VaultSecret{
							"test":            {},
							v1alpha1.GitToken: {},
						},
					},
				},
			}
			err = p.ValidateCreate()
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(ContainSubstring("Invalid value: test"))
		})
	})
	Context("various", func() {
		It("returns the stack string", func() {
			p := &v1alpha1.Play{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
			}
			_ = p.ValidateDelete()
			s := v1alpha1.Stack{
				Language: "js",
				Package:  "npm",
			}
			Expect(s.String()).To(Equal("js/npm"))
		})
		It("sets default", func() {
			p := &v1alpha1.Play{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: v1alpha1.PlaySpec{
					Scope: v1alpha1.Scope{
						Name: "",
					},
				},
			}
			p.Default()
			Expect(p.Spec.Scope.Name).To(Equal("default"))
		})
		It("success", func() {
			var err error
			p := &v1alpha1.Play{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: v1alpha1.PlaySpec{
					Scope: v1alpha1.Scope{
						Name: "",
					},
					Environment: "develop",
					ProjectID:   1,
					PipelineID:  1,
					RepoURL:     "https://github.com",
					Commit: v1alpha1.Commit{
						SHA: "sha",
						Ref: "main",
					},
					Tasks: []map[v1alpha1.TaskType]v1alpha1.Task{
						{
							v1alpha1.Deploy: v1alpha1.Task{},
						},
					},
				},
			}
			old := &v1alpha1.Play{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: v1alpha1.PlaySpec{
					Scope: v1alpha1.Scope{
						Name: "",
					},
					Environment: "develop",
					ProjectID:   1,
					PipelineID:  1,
					RepoURL:     "https://github.com",
					Commit: v1alpha1.Commit{
						SHA: "sha",
						Ref: "main",
					},
					Tasks: []map[v1alpha1.TaskType]v1alpha1.Task{
						{
							v1alpha1.Deploy: v1alpha1.Task{},
						},
					},
				},
			}
			err = p.ValidateUpdate(old)
			Expect(err).To(Succeed())
		})

		It("fails on update", func() {
			var err error
			p := &v1alpha1.Play{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: v1alpha1.PlaySpec{
					Scope: v1alpha1.Scope{
						Name: "",
					},
					ProjectID:  1,
					PipelineID: 1,
					RepoURL:    "http://{}",
					Domain:     "ee_@",
					Tasks: []map[v1alpha1.TaskType]v1alpha1.Task{
						{
							"test": v1alpha1.Task{},
						},
						{
							v1alpha1.Deploy: v1alpha1.Task{},
						},
					},
				},
			}
			old := &v1alpha1.Play{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: v1alpha1.PlaySpec{
					Scope: v1alpha1.Scope{
						Name: "",
					},
					RepoURL: "https://test.fr",
					Domain:  "test.fr",
					Tasks: []map[v1alpha1.TaskType]v1alpha1.Task{
						{
							v1alpha1.Deploy: v1alpha1.Task{},
						},
					},
				},
			}
			err = p.ValidateUpdate(old)
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(ContainSubstring("Invalid value: test"))
		})
		It("checks", func() {
			steps := v1alpha1.Steps{
				{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							v1alpha1.AnnotationOrder: "1",
						},
					},
					Step: v1alpha1.StepSpec{
						Step: tkn.Step{
							Container: corev1.Container{
								Image: "step1",
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							v1alpha1.AnnotationOrder: "2",
						},
					},
					Step: v1alpha1.StepSpec{
						Step: tkn.Step{
							Container: corev1.Container{
								Image: "step2",
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							v1alpha1.AnnotationOrder: "a",
						},
					},
					Step: v1alpha1.StepSpec{
						Step: tkn.Step{
							Container: corev1.Container{
								Image: "step2",
							},
						},
					},
				},
			}
			Expect(steps.Len()).To(Equal(3))
			Expect(steps.Less(0, 1)).To(Equal(true))
			Expect(steps.Less(1, 0)).To(Equal(false))
			Expect(steps.Less(1, 2)).To(Equal(false))
			Expect(steps.Less(2, 1)).To(Equal(true))
			steps.Swap(0, 1)
			Expect(steps[0].Step.Step.Image).To(Equal("step2"))
		})
	})
})
