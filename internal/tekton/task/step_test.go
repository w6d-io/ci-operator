/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 12/02/2021
*/

package task_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/k8s/secrets"
	"github.com/w6d-io/ci-operator/internal/tekton/task"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

var _ = Describe("Task", func() {
	Context("validate GetGeneric behaviour", func() {
		var (
			s      *task.Step
			logger = ctrl.Log.WithName("test")
		)
		BeforeEach(func() {
			s = &task.Step{
				Client: k8sClient,
			}
		})
		It("no in the same namespace", func() {
			config.SetNamespace("test")
			steps := ci.Steps{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
					},
				},
			}
			gets := s.GetGenericSteps(logger, steps)
			Expect(len(gets)).To(Equal(0))
		})
		It("no annotation kind", func() {
			config.SetNamespace("test")
			steps := ci.Steps{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "test",
					},
				},
			}
			gets := s.GetGenericSteps(logger, steps)
			Expect(len(gets)).To(Equal(0))
		})
		It("no the same task", func() {
			config.SetNamespace("test")
			s.TaskType = ci.UnitTests
			steps := ci.Steps{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "test",
						Annotations: map[string]string{
							ci.AnnotationKind: "generic",
							ci.AnnotationTask: "build",
						},
					},
				},
			}
			gets := s.GetGenericSteps(logger, steps)
			Expect(len(gets)).To(Equal(0))
		})
		It("get a step", func() {
			config.SetNamespace("test")
			s.TaskType = ci.Build
			steps := ci.Steps{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "test",
						Annotations: map[string]string{
							ci.AnnotationKind: "generic",
							ci.AnnotationTask: ci.Build.String(),
						},
					},
				},
			}
			gets := s.GetGenericSteps(logger, steps)
			Expect(len(gets)).To(Equal(1))
		})
	})
	Context("Step methods", func() {
		It("deal with FilteredSteps", func() {
			By("set namespace")
			config.SetNamespace("default")

			By("set task step")
			s := &task.Step{
				Index: 0,
				PlaySpec: ci.PlaySpec{
					Stack: ci.Stack{
						Language: "bash",
						Package:  "test",
					},
					Tasks: []map[ci.TaskType]ci.Task{
						{
							ci.UnitTests: ci.Task{
								Script: ci.Script{
									"echo", "toto",
								},
							},
						},
					},
				},
				TaskType: ci.E2ETests,
			}

			By("set ci step  with unmatched namespace")
			steps := ci.Steps{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "p6e-cx-20",
						Annotations: map[string]string{
							ci.AnnotationPackage: "null",
						},
					},
				},
			}
			Expect(len(s.FilteredSteps(ctrl.Log, steps, false))).To(Equal(0))

			By("set namespace")
			config.SetNamespace("p6e-cx-20")

			By("package not match")
			Expect(len(s.FilteredSteps(ctrl.Log, steps, true))).To(Equal(0))

			By("package match")
			s.TaskType = ci.UnitTests
			Expect(len(s.FilteredSteps(ctrl.Log, steps, true))).To(Equal(0))

			By("task unmatched")
			s.TaskType = ci.Build
			Expect(len(s.FilteredSteps(ctrl.Log, steps, false))).To(Equal(0))

			By("language unmatched")
			s.TaskType = ci.UnitTests
			steps[0].Annotations = map[string]string{
				ci.AnnotationTask:     ci.UnitTests.String(),
				ci.AnnotationPackage:  "test",
				ci.AnnotationLanguage: "none",
			}
			Expect(len(s.FilteredSteps(ctrl.Log, steps, false))).To(Equal(0))

			By("return a step")
			By("language unmatched")
			s.TaskType = ci.UnitTests
			steps[0].Annotations = map[string]string{
				ci.AnnotationTask:     ci.UnitTests.String(),
				ci.AnnotationPackage:  "test",
				ci.AnnotationLanguage: "bash",
			}
			Expect(len(s.FilteredSteps(ctrl.Log, steps, false))).To(Equal(1))
		})
		It("deals with GetSteps", func() {
			var err error
			By("build step")
			s := &task.Step{
				Index: 0,
				PlaySpec: ci.PlaySpec{
					Stack: ci.Stack{
						Language: "bash",
						Package:  "test",
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
					Vault: &ci.Vault{
						Secrets: map[ci.SecretKind]ci.VaultSecret{
							secrets.KubeConfigKey: {},
						},
					},
				},
				TaskType: ci.E2ETests,
				Client:   k8sClient,
			}

			By("Create namespace")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "p6e-cx-21",
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(Succeed())

			By("Create step")
			step := &ci.Step{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "step-unit-test-1",
					Namespace: "p6e-cx-21",
					Annotations: map[string]string{
						ci.AnnotationPackage:  "test",
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

			By("set task type")
			s.TaskType = ci.Build

			By("Set config")
			Expect(config.New("testdata/config.yaml")).To(Succeed())

			By("Set namespace")
			config.SetNamespace("p6e-cx-21")

			_, _, _, err = s.GetSteps(ctx, ctrl.Log)
			Expect(err).To(Succeed())

			By("set fake client")
			s.Client = fake.NewClientBuilder().Build()

			_, _, _, err = s.GetSteps(ctx, ctrl.Log)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no kind is registered for the type v1alpha1.StepList"))

			By("test get params")
			params := []ci.ParamSpec{
				{
					ParamSpec: tkn.ParamSpec{
						Name:        "test1",
						Type:        "string",
						Description: "unit test get params",
						Default: &tkn.ArrayOrString{
							StringVal: "no default",
						},
					},
				},
			}
			Expect(len(s.GetParams(params))).To(Equal(1))
		})
	})
})
