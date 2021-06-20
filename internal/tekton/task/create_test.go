/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 30/04/2021
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
	ctrl "sigs.k8s.io/controller-runtime"
)

var _ = Describe("Task", func() {
	Context("Parse", func() {
		BeforeEach(func() {
		})
		AfterEach(func() {
		})
		It("", func() {
			var err error

			By("create task")
			t := task.Task{
				Client: k8sClient,
				Play: &ci.Play{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "play-parse-1",
						Namespace: "p6e-cx-40",
						UID:       "11d3a4b9-a40b-4a63-9940-44a3ca6bc254",
					},
					Spec: ci.PlaySpec{
						Stack: ci.Stack{
							Language: "bash",
							Package:  "custom",
						},
						Tasks: []map[ci.TaskType]ci.Task{
							{
								ci.Build: ci.Task{
									Script: ci.Script{
										"echo", "toto",
									},
								},
							},
						},
					},
				},
				Scheme: scheme,
			}
			By("Set namespace")
			config.SetNamespace("p6e-cx-40")
			err = t.Parse(ctx, ctrl.Log)
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(Equal("no step found for build"))

			By("Create namespace")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "p6e-cx-40",
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(Succeed())

			By("Create build step")
			step := &ci.Step{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "step-build-1",
					Namespace: "p6e-cx-40",
					Annotations: map[string]string{
						ci.AnnotationPackage:  "custom",
						ci.AnnotationLanguage: "bash",
						ci.AnnotationTask:     ci.Build.String(),
						ci.AnnotationOrder:    "0",
					},
				},
				Step: ci.StepSpec{
					Step: tkn.Step{
						Script: "echo test",
					},
				},
			}
			Expect(k8sClient.Create(ctx, step)).To(Succeed())
			err = t.Parse(ctx, ctrl.Log)
			Expect(err).To(Succeed())

			// UnitTest
			By("Create namespace")
			ns = &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "p6e-cx-42",
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(Succeed())
			t.Play.Spec.Tasks = []map[ci.TaskType]ci.Task{
				{
					ci.UnitTests: ci.Task{
						Script: ci.Script{
							"echo", "toto",
						},
					},
				},
			}
			t.Play.Spec.ProjectID = 42
			t.Play.Namespace = "p6e-cx-42"
			err = t.Parse(ctx, ctrl.Log)
			Expect(err).ToNot(Succeed())
			By("Set namespace")
			config.SetNamespace("p6e-cx-42")
			Expect(err.Error()).To(Equal("no step found for unit-tests"))
			By("Create unit test step")
			step = &ci.Step{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "step-unit-test-1",
					Namespace: "p6e-cx-42",
					Annotations: map[string]string{
						ci.AnnotationPackage:  "custom",
						ci.AnnotationLanguage: "bash",
						ci.AnnotationTask:     ci.UnitTests.String(),
						ci.AnnotationOrder:    "0",
					},
				},
				Step: ci.StepSpec{
					Step: tkn.Step{
						Script: "echo test",
					},
				},
			}
			Expect(k8sClient.Create(ctx, step)).To(Succeed())
			err = t.Parse(ctx, ctrl.Log)
			Expect(err).To(Succeed())

			// Integration Test
			By("Create namespace")
			ns = &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "p6e-cx-43",
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(Succeed())
			t.Play.Spec.Tasks = []map[ci.TaskType]ci.Task{
				{
					ci.IntegrationTests: ci.Task{
						Script: ci.Script{
							"echo", "toto",
						},
					},
				},
			}
			t.Play.Spec.ProjectID = 43
			t.Play.Namespace = "p6e-cx-43"
			err = t.Parse(ctx, ctrl.Log)
			Expect(err).ToNot(Succeed())
			By("Set namespace")
			config.SetNamespace("p6e-cx-43")
			Expect(err.Error()).To(Equal("no step found for integration-tests"))
			By("Create integration test step")
			step = &ci.Step{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "step-integration-test-1",
					Namespace: "p6e-cx-43",
					Annotations: map[string]string{
						ci.AnnotationPackage:  "custom",
						ci.AnnotationLanguage: "bash",
						ci.AnnotationTask:     ci.IntegrationTests.String(),
						ci.AnnotationOrder:    "0",
					},
				},
				Step: ci.StepSpec{
					Step: tkn.Step{
						Script: "echo test",
					},
				},
			}
			Expect(k8sClient.Create(ctx, step)).To(Succeed())
			err = t.Parse(ctx, ctrl.Log)
			Expect(err).To(Succeed())

			// Deploy
			By("Create namespace")
			ns = &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "p6e-cx-44",
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(Succeed())
			t.Play.Spec.Tasks = []map[ci.TaskType]ci.Task{
				{
					ci.Deploy: ci.Task{
						Script: ci.Script{
							"echo", "toto",
						},
					},
				},
			}
			t.Play.Spec.ProjectID = 44
			t.Play.Namespace = "p6e-cx-44"
			err = t.Parse(ctx, ctrl.Log)
			Expect(err).ToNot(Succeed())
			By("Set namespace")
			config.SetNamespace("p6e-cx-44")
			Expect(err.Error()).To(Equal("no step found for deploy"))
			By("Create deploy test step")
			step = &ci.Step{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "step-deploy-test-1",
					Namespace: "p6e-cx-44",
					Annotations: map[string]string{
						ci.AnnotationPackage:  "custom",
						ci.AnnotationLanguage: "bash",
						ci.AnnotationTask:     ci.Deploy.String(),
						ci.AnnotationOrder:    "0",
					},
				},
				Step: ci.StepSpec{
					Step: tkn.Step{
						Script: "echo test",
					},
				},
			}
			Expect(k8sClient.Create(ctx, step)).To(Succeed())
			err = t.Parse(ctx, ctrl.Log)
			Expect(err).To(Succeed())

			// Clean
			By("Create namespace")
			ns = &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "p6e-cx-45",
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(Succeed())
			t.Play.Spec.Tasks = []map[ci.TaskType]ci.Task{
				{
					ci.Clean: ci.Task{
						Script: ci.Script{
							"echo", "toto",
						},
					},
				},
			}
			t.Play.Spec.ProjectID = 45
			t.Play.Namespace = "p6e-cx-45"
			err = t.Parse(ctx, ctrl.Log)
			Expect(err).ToNot(Succeed())
			By("Set namespace")
			config.SetNamespace("p6e-cx-45")
			Expect(err.Error()).To(Equal("no step found for clean"))
			By("Create clean test step")
			step = &ci.Step{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "step-clean-test-1",
					Namespace: "p6e-cx-45",
					Annotations: map[string]string{
						ci.AnnotationPackage:  "custom",
						ci.AnnotationLanguage: "bash",
						ci.AnnotationTask:     ci.Clean.String(),
						ci.AnnotationOrder:    "0",
					},
				},
				Step: ci.StepSpec{
					Step: tkn.Step{
						Script: "echo test",
					},
				},
			}
			Expect(k8sClient.Create(ctx, step)).To(Succeed())
			err = t.Parse(ctx, ctrl.Log)
			Expect(err).To(Succeed())

			// E2ETest
			By("Create namespace")
			ns = &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "p6e-cx-46",
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(Succeed())
			t.Play.Spec.Tasks = []map[ci.TaskType]ci.Task{
				{
					ci.E2ETests: ci.Task{
						Script: ci.Script{
							"echo", "toto",
						},
					},
				},
			}
			t.Play.Spec.ProjectID = 46
			t.Play.Namespace = "p6e-cx-46"
			err = t.Parse(ctx, ctrl.Log)
			Expect(err).ToNot(Succeed())
			By("Set namespace")
			config.SetNamespace("p6e-cx-46")
			Expect(err.Error()).To(Equal("no step found for e2e-tests"))
			By("Create 2e2 test step")
			step = &ci.Step{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "step-e2e-test-1",
					Namespace: "p6e-cx-46",
					Annotations: map[string]string{
						ci.AnnotationPackage:  "custom",
						ci.AnnotationLanguage: "bash",
						ci.AnnotationTask:     ci.E2ETests.String(),
						ci.AnnotationOrder:    "0",
					},
				},
				Step: ci.StepSpec{
					Step: tkn.Step{
						Script: "echo test",
					},
				},
			}
			Expect(k8sClient.Create(ctx, step)).To(Succeed())
			err = t.Parse(ctx, ctrl.Log)
			Expect(err).To(Succeed())

			// Generic
			By("Create namespace")
			ns = &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "p6e-cx-47",
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(Succeed())
			t.Play.Spec.Tasks = []map[ci.TaskType]ci.Task{
				{
					"git-leaks": ci.Task{
						Script: ci.Script{
							"echo", "toto",
						},
					},
				},
			}
			t.Play.Spec.ProjectID = 47
			t.Play.Namespace = "p6e-cx-47"
			t.Params = map[string][]ci.ParamSpec{
				"test": {
					ci.ParamSpec{
						ParamSpec: tkn.ParamSpec{
							Name: "unit-test",
						},
					},
				},
			}
			err = t.Parse(ctx, ctrl.Log)
			Expect(err).ToNot(Succeed())

			By("Set namespace")
			config.SetNamespace("p6e-cx-47")
			Expect(err.Error()).To(Equal("no step found for git-leaks"))
			By("Create git-leaks test step")
			step = &ci.Step{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "step-git-leaks",
					Namespace:    "p6e-cx-47",
					Annotations: map[string]string{
						ci.AnnotationLanguage: "bash",
						ci.AnnotationPackage: "custom",
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
			Expect(k8sClient.Create(ctx, step)).To(Succeed())
			err = t.Parse(ctx, ctrl.Log)
			Expect(err).To(Succeed())

		})
	})
})
