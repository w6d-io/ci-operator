/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 20/05/2021
*/

package v1alpha1_test

import (
	. "github.com/onsi/ginkgo"
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/w6d-io/ci-operator/api/v1alpha1"
)

//goland:noinspection GoNilness,GoNilness,GoNilness,GoNilness,GoNilness,GoNilness,GoNilness,GoNilness,GoNilness,GoNilness,GoNilness,GoNilness,GoNilness,GoNilness,GoNilness,GoNilness,GoNilness,GoNilness,GoNilness,GoNilness,GoNilness,GoNilness,GoNilness,GoNilness,GoNilness,GoNilness,GoNilness
var _ = Describe("DeepCopy", func() {
	Context("", func() {
		BeforeEach(func() {
		})
		AfterEach(func() {
		})
		It("Commit", func() {
			By("DeepCopyInto")
			in := &v1alpha1.Commit{
				SHA: "test",
			}
			out := &v1alpha1.Commit{}
			in.DeepCopyInto(out)

			By("DeepCopy")
			_ = in.DeepCopy()

			By("return nil")
			var null *v1alpha1.Commit
			//goland:noinspection GoNilness
			_ = null.DeepCopy()
		})
		It("Docker", func() {
			By("DeepCopyInto")
			in := &v1alpha1.Docker{
				Context: "test",
			}
			out := &v1alpha1.Docker{}
			in.DeepCopyInto(out)

			By("DeepCopy")
			_ = in.DeepCopy()

			By("return nil")
			var null *v1alpha1.Docker
			_ = null.DeepCopy()
		})
		It("LimitCi", func() {
			By("DeepCopyInto")
			in := &v1alpha1.LimitCi{
				Spec: v1alpha1.LimitCiSpec{Concurrent: 1},
			}
			out := &v1alpha1.LimitCi{}
			in.DeepCopyInto(out)

			in.DeepCopyObject()

			By("DeepCopy")
			_ = in.DeepCopy()

			By("return nil")
			var null *v1alpha1.LimitCi
			_ = null.DeepCopy()

			_ = null.DeepCopyObject()

		})
		It("LimitCiList", func() {
			By("DeepCopyInto")
			in := &v1alpha1.LimitCiList{
				Items: []v1alpha1.LimitCi{
					{
						Spec: v1alpha1.LimitCiSpec{
							Concurrent: 1,
						},
					},
				},
			}
			out := &v1alpha1.LimitCiList{}
			in.DeepCopyInto(out)

			in.DeepCopyObject()

			By("DeepCopy")
			_ = in.DeepCopy()

			By("return nil")
			var null *v1alpha1.LimitCiList
			_ = null.DeepCopy()

			_ = null.DeepCopyObject()

		})
		It("LimitCiSpec", func() {
			By("DeepCopyInto")
			in := &v1alpha1.LimitCiSpec{}
			out := &v1alpha1.LimitCiSpec{}
			in.DeepCopyInto(out)

			By("DeepCopy")
			_ = in.DeepCopy()

			By("return nil")
			var null *v1alpha1.LimitCiSpec
			_ = null.DeepCopy()

		})
		It("LimitCiStatus", func() {
			By("DeepCopyInto")
			in := &v1alpha1.LimitCiStatus{}
			out := &v1alpha1.LimitCiStatus{}
			in.DeepCopyInto(out)

			By("DeepCopy")
			_ = in.DeepCopy()

			By("return nil")
			var null *v1alpha1.LimitCiStatus
			_ = null.DeepCopy()

		})
		It("ParamSpec", func() {
			By("DeepCopyInto")
			in := &v1alpha1.ParamSpec{}
			out := &v1alpha1.ParamSpec{}
			in.DeepCopyInto(out)

			By("DeepCopy")
			_ = in.DeepCopy()

			By("return nil")
			var null *v1alpha1.ParamSpec
			_ = null.DeepCopy()

		})
		It("Play", func() {
			By("DeepCopyInto")
			in := &v1alpha1.Play{
				Spec: v1alpha1.PlaySpec{},
			}
			out := &v1alpha1.Play{}
			in.DeepCopyInto(out)

			in.DeepCopyObject()

			By("DeepCopy")
			_ = in.DeepCopy()

			By("return nil")
			var null *v1alpha1.Play
			_ = null.DeepCopy()

			_ = null.DeepCopyObject()

		})
		It("PlayList", func() {
			By("DeepCopyInto")
			in := &v1alpha1.PlayList{
				Items: []v1alpha1.Play{
					{
						Spec: v1alpha1.PlaySpec{
							Tasks: []map[v1alpha1.TaskType]v1alpha1.Task{
								{"task": v1alpha1.Task{}},
							},
							Secret: map[v1alpha1.SecretKind]string{
								"secret": "secret",
							},
							Vault: &v1alpha1.Vault{
								Token: "token",
							},
						},
						Status: v1alpha1.PlayStatus{
							PipelineRunName: "pr-name",
						},
					},
				},
			}
			out := &v1alpha1.PlayList{}
			in.DeepCopyInto(out)

			in.DeepCopyObject()

			By("DeepCopy")
			_ = in.DeepCopy()

			By("return nil")
			var null *v1alpha1.PlayList
			_ = null.DeepCopy()

			_ = null.DeepCopyObject()
		})
		It("PlaySpec", func() {
			spec := &v1alpha1.PlaySpec{
				Tasks: []map[v1alpha1.TaskType]v1alpha1.Task{
					{"task": v1alpha1.Task{}},
				},
				Secret: map[v1alpha1.SecretKind]string{
					"secret": "secret",
				},
				Vault: &v1alpha1.Vault{
					Token: "token",
				},
			}

			By("DeepCopy")
			_ = spec.DeepCopy()

			var out *v1alpha1.PlaySpec
			out.DeepCopy()
		})
	})
	It("PlayStatus", func() {
		By("DeepCopyInto")
		in := &v1alpha1.PlayStatus{
			PipelineRunName: "pr-name",
		}
		out := &v1alpha1.PlayStatus{}
		in.DeepCopyInto(out)

		By("DeepCopy")
		_ = in.DeepCopy()

		By("return nil")
		var null *v1alpha1.PlayStatus
		_ = null.DeepCopy()

	})
	It("Scope", func() {
		By("DeepCopyInto")
		in := &v1alpha1.Scope{}
		out := &v1alpha1.Scope{}
		in.DeepCopyInto(out)

		By("DeepCopy")
		_ = in.DeepCopy()

		By("return nil")
		var null *v1alpha1.Scope
		_ = null.DeepCopy()

	})
	It("Script", func() {
		By("DeepCopyInto")
		in := &v1alpha1.Script{}
		out := &v1alpha1.Script{}
		in.DeepCopyInto(out)

		By("DeepCopy")
		_ = in.DeepCopy()

		By("return nil")
		var null v1alpha1.Script
		_ = null.DeepCopy()

	})
	It("Stack", func() {
		By("DeepCopyInto")
		in := &v1alpha1.Stack{}
		out := &v1alpha1.Stack{}
		in.DeepCopyInto(out)

		By("DeepCopy")
		_ = in.DeepCopy()

		By("return nil")
		var null *v1alpha1.Stack
		_ = null.DeepCopy()

	})
	It("Step", func() {
		By("DeepCopyInto")
		in := &v1alpha1.Step{
			Params: []v1alpha1.ParamSpec{
				{},
			},
		}
		out := &v1alpha1.Step{}

		in.DeepCopyInto(out)

		By("DeepCopy")
		_ = in.DeepCopy()

		By("return nil")
		var null *v1alpha1.Step
		_ = null.DeepCopy()
		_ = null.DeepCopyObject()

		By("DeepCopyObject")
		_ = in.DeepCopyObject()
	})
	It("StepList", func() {
		By("DeepCopyInto")
		in := &v1alpha1.StepList{
			Items: []v1alpha1.Step{
				{},
			},
		}
		out := &v1alpha1.StepList{}
		in.DeepCopyInto(out)

		in.DeepCopyObject()

		By("DeepCopy")
		_ = in.DeepCopy()

		By("return nil")
		var null *v1alpha1.StepList
		_ = null.DeepCopy()

		_ = null.DeepCopyObject()
	})
	It("StepSpec", func() {
		var out *v1alpha1.StepSpec
		_ = out.DeepCopy()

		in := &v1alpha1.StepSpec{
			Step: tkn.Step{
				Script: "",
			},
		}
		_ = in.DeepCopy()

	})
	It("StepStatus", func() {
		By("DeepCopyInto")
		in := &v1alpha1.StepStatus{}
		out := &v1alpha1.StepStatus{}
		in.DeepCopyInto(out)

		By("DeepCopy")
		_ = in.DeepCopy()

		By("return nil")
		var null *v1alpha1.StepStatus
		_ = null.DeepCopy()

	})
	It("Steps", func() {
		By("DeepCopyInto")
		in := &v1alpha1.Steps{
			{
				Step: v1alpha1.StepSpec{
					Step: tkn.Step{
						Script: "",
					},
				},
			},
		}
		out := &v1alpha1.Steps{}
		in.DeepCopyInto(out)

		By("DeepCopy")
		_ = in.DeepCopy()

		By("return nil")
		var null v1alpha1.Steps
		_ = null.DeepCopy()
	})
	It("Task", func() {
		in := &v1alpha1.Task{
			Script: v1alpha1.Script{
				"test", "test1",
			},
			Arguments: []string{
				"arg1", "args2",
			},
			Variables: map[string]string{
				"NAME": "value",
			},
			Annotations: map[string]string{
				"key": "value",
			},
		}
		out := &v1alpha1.Task{}
		By("DeepCopyInto")
		in.DeepCopyInto(out)

		var null *v1alpha1.Task
		_ = null.DeepCopy()
	})
	It("Vault", func() {
		in := &v1alpha1.Vault{
			Token: "token",
			Secrets: map[v1alpha1.SecretKind]v1alpha1.VaultSecret{
				v1alpha1.KubeConfig: {
					Path: "/root",
				},
			},
		}
		out := &v1alpha1.Vault{}
		in.DeepCopyInto(out)
		_ = in.DeepCopy()
		var null *v1alpha1.Vault
		null.DeepCopy()

	})
	It("VaultSecret", func() {
		in := &v1alpha1.VaultSecret{
			Path: "/root",
		}
		out := &v1alpha1.VaultSecret{}
		in.DeepCopyInto(out)
		_ = in.DeepCopy()
		var null *v1alpha1.VaultSecret
		null.DeepCopy()

	})
})
