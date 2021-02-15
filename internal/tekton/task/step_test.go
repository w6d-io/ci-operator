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
	"github.com/w6d-io/ci-operator/internal/config"

	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/w6d-io/ci-operator/internal/tekton/task"
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
})
