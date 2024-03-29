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
Created on 04/03/2021
*/
package webhook_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	ci "github.com/w6d-io/ci-operator/api/v1alpha1"

	"github.com/w6d-io/ci-operator/pkg/webhook"
)

var _ = Describe("Webhook", func() {
	Context("limit ci", func() {
		It("get payload", func() {
			p := &ci.Play{
				Spec: ci.PlaySpec{
					ProjectID:  1,
					PipelineID: 1,
				},
				Status: ci.PlayStatus{
					State: ci.Succeeded,
				},
			}
			l := ci.LimitCi{
				Spec: ci.LimitCiSpec{
					Concurrent: 2,
				},
			}
			m := "you reach the concurrent limit"
			payload := webhook.GetLimitPayload(p, l, m)
			Expect(payload).ToNot(BeNil())
		})
	})
})
