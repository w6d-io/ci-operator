/*
Copyright 2020 WILDCARD SA.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
Created on 02/03/2021
*/
package pipelinerun_test

import (
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/w6d-io/ci-operator/internal"
	"github.com/w6d-io/ci-operator/internal/tekton/pipelinerun"

	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("generic in pipeline run", func() {
	Context("setting", func() {
		It("does", func() {
			p := pipelinerun.PipelineRun{
				WorkFlowStruct: internal.WorkFlowStruct{
					Play: &ci.Play{
						Spec: ci.PlaySpec{
							Name:       "test",
							ProjectID:  1,
							PipelineID: 1,
							Stack: ci.Stack{
								Language: "test",
							},
							Commit: ci.Commit{
								SHA: "test_test_test",
								Ref: "test",
							},
							Tasks: []map[ci.TaskType]ci.Task{
								{
									"test": ci.Task{
										Arguments: []string{
											"TEST", "Test",
										},
										Image: "test/test:test",
										Script: []string{
											"echo",
											"test",
										},
									},
								},
							},
						},
					},
				},
				GenericParams: map[string][]ci.ParamSpec{
					"test": {
						ci.ParamSpec{
							ParamSpec: tkn.ParamSpec{
								Name: "test",
							},
							Value: "value",
						},
					},
				},
			}
			err := p.SetGeneric(0, "test", ctrl.Log)
			Expect(err).To(Succeed())
		})
	})
})
