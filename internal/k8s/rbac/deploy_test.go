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
package rbac_test

import (
	. "github.com/onsi/ginkgo"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/k8s/rbac"
)

var _ = Describe("RBAC", func() {
	Context("", func() {
		BeforeEach(func() {
		})
		AfterEach(func() {
		})
		config.New("testdata/config.yaml")
		It("", func() {
			p := &ci.Play{
				Spec: ci.PlaySpec{
					ProjectID:  1,
					PipelineID: 1,
				},
			}
			_ = rbac.GetRoleBinding(p)
		})
		It("", func() {
			p := &ci.Play{
				Spec: ci.PlaySpec{
					ProjectID:  1,
					PipelineID: 1,
				},
			}
			_ = rbac.GetSubject(p)
		})
	})
})
