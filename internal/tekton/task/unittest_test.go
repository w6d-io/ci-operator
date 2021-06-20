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

package task_test

import (
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/k8s/secrets"
	"github.com/w6d-io/ci-operator/internal/tekton/task"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Task", func() {
	Context("UnitTests", func() {
		It("execute UnitTests", func() {
			var err error
			By("Create task")
			t := &task.Task{
				Index:  0,
				Client: k8sClient,
				Play: &ci.Play{
					Spec: ci.PlaySpec{
						Stack: ci.Stack{
							Language: "bash",
							Package:  "test",
						},
						Tasks: []map[ci.TaskType]ci.Task{
							{
								ci.UnitTests: ci.Task{
									Variables: fields.Set{
										"TEST": "OK",
									},
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
				},
			}

			By("load config")
			Expect(config.New("testdata/config.yaml")).To(Succeed())

			By("Set namespace")
			config.SetNamespace("p6e-cx-30")

			By("Create namespace")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "p6e-cx-30",
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(Succeed())

			By("failed by no step")
			err = t.UnitTest(ctx, ctrl.Log)
			Expect(err).ToNot(Succeed())

			By("Create step")
			step := &ci.Step{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "step-unit-test-1",
					Namespace: "p6e-cx-30",
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

			Expect(t.UnitTest(ctx, ctrl.Log)).To(Succeed())

			By("set index overflow")
			t.Index = 3
			err = t.UnitTest(ctx, ctrl.Log)
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(Equal("no such task"))
		})
		It("Execute Create", func() {
			var err error
			By("build UnitTestsTask")
			u := &task.UnitTestTask{
				Meta: task.Meta{
					Steps: []tkn.Step{
						{
							Script: "echo test",
						},
					},
					Play: &ci.Play{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "play-test-31-1",
							Namespace: "p6e-cx-31",
							UID:       "cf3e4131-6b00-410e-9d3d-774312a57bce",
						},
						Spec: ci.PlaySpec{
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
							Vault: &ci.Vault{
								Secrets: map[ci.SecretKind]ci.VaultSecret{
									secrets.KubeConfigKey: {},
								},
							},
						},
					},
					Scheme: scheme,
				},
			}

			By("Set config")
			Expect(config.New("testdata/config.yaml")).To(Succeed())

			By("Failed due to kind missing")
			u.Scheme = runtime.NewScheme()
			err = u.Create(ctx, k8sClient, ctrl.Log)
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(ContainSubstring("no kind is registered for the type v1alpha1.Play"))

			By("set the right namespace")
			u.Scheme = scheme
			u.Play.Spec.PipelineID = 1
			u.Play.Spec.ProjectID = 31
			err = u.Create(ctx, k8sClient, ctrl.Log)
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(Equal(`namespaces "p6e-cx-31" not found`))

			By("create namespace")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "p6e-cx-31",
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(Succeed())

			Expect(u.Create(ctx, k8sClient, ctrl.Log)).To(Succeed())

		})
	})
})
