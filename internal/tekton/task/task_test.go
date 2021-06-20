/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 19/06/2021
*/

package task_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/tekton/task"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("Task", func() {
	Context("Generic", func() {
		BeforeEach(func() {
		})
		AfterEach(func() {
		})
		It("Create a generic task", func() {
			var err error
			By("Build a task")
			t := &task.Task{
				Index: 0,
				Play: &ci.Play{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "play-unit-test-generic",
						Namespace: "p6e-cx-80",
						UID:       "aaa-bbb-ccc-ddd",
					},
					Spec: ci.PlaySpec{
						Stack: ci.Stack{
							Language: "python",
							Package:  "pip",
						},
						ProjectID:  80,
						PipelineID: 1,
						Commit: ci.Commit{
							SHA: "aaa-bbb-ccc-ddd-eee-fff-ggg",
						},
						Tasks: []map[ci.TaskType]ci.Task{
							{
								"git-leaks": ci.Task{
									Script: ci.Script{
										"echo", "unit-test",
									},
									Variables: fields.Set{
										"VAR": "content",
									},
								},
							},
						},
					},
				},
				Scheme: scheme,
				Params: map[string][]ci.ParamSpec{
					"test": {
						{
							ParamSpec: tkn.ParamSpec{
								Name: "test1",
								Type: "string",
								Default: &tkn.ArrayOrString{
									StringVal: "value1",
								},
							},
						},
					},
					"git-leaks": {
						{
							ParamSpec: tkn.ParamSpec{
								Name: "test1",
								Type: "string",
								Default: &tkn.ArrayOrString{
									StringVal: "value1",
								},
							},
						},
					},
				},
			}

			By("Call with fake client")
			t.Client = fake.NewClientBuilder().Build()
			err = t.Generic(ctx, "git-leaks", ctrl.Log)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no kind is registered for the type v1alpha1.StepList"))

			By("Call Generic with no such step")
			t.Client = k8sClient
			err = t.Generic(ctx, "git-leaks", ctrl.Log)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("no step found for git-leaks"))

			By("Create namespace")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "unit-test",
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(Succeed())
			config.SetNamespace(ns.GetName())

			By("Create step")
			step := &ci.Step{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "step-git-leaks",
					Namespace:    ns.GetName(),
					Annotations: map[string]string{
						ci.AnnotationLanguage: "python",
						ci.AnnotationPackage: "pip",
						ci.AnnotationOrder: "0",
						ci.AnnotationTask:  "git-leaks",
					},
				},
				Params: []ci.ParamSpec{
					{
						ParamSpec: tkn.ParamSpec{},
					},
				},
				Step: ci.StepSpec{
					Step: tkn.Step{
						Container: corev1.Container{
							Name:  "git-leaks",
							Image: "w6dio/docker-gitleaks:v0.0.6",
						},
						Script: "echo git-leaks",
					},
				},
			}
			Expect(k8sClient.Create(ctx, step))
			Expect(t.Generic(ctx, "git-leaks", ctrl.Log)).To(Succeed())

			By("Call create")
			Expect(len(t.Creates)).To(Equal(1))

			f := t.Creates[0]
			err = f(ctx, k8sClient, ctrl.Log)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`namespaces "p6e-cx-80" not found`))

			t.Play.Namespace = "default"
			f = t.Creates[0]
			err = f(ctx, k8sClient, ctrl.Log)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`cross-namespace owner references are disallowed, owner's namespace default, obj's namespace p6e-cx-80`))

			t.Play.Namespace = ns.GetName()
			t.Play.Spec.PipelineNamespace = ns.GetName()
			f = t.Creates[0]
			Expect(f(ctx, k8sClient, ctrl.Log)).To(Succeed())

		})
	})
})
