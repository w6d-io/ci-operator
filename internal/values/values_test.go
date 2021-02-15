/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 15/02/2021
*/

package values_test

import (
	"bytes"
	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/values"
)

var playSpec = `
commit:
  message: no message
  ref: main
  sha: 3d5531613d5531613d5531613d5531613d553161
environment: staging
name: test
pipeline_id: 1
project_id: 1
docker_url: reg.example.com/group/repo:test
repo_url: https://example.com/group/repo
scope: {}
secret:
  .dockerconfigjson: '{"auths":{"reg.example.com":{"auth":"dGVzdDp0ZXN0Cg=="}}}'
  git_token: 3d5531613d5531613d5531613d5531613d553161
tasks:
- deploy:
    docker: {}
    variables:
      TEST: test
dependencies:
  mongodb:
    variables:
      HOST: "$DATABASE_HOST"
`

var _ = Describe("Values", func() {
	Context("Template", func() {
		var (
			p *ci.Play
		)
		BeforeEach(func() {
			spec := ci.PlaySpec{}
			err := yaml.Unmarshal([]byte(playSpec), &spec)
			Expect(err).To(Succeed())
			p = &ci.Play{
				Spec: spec,
			}
		})
		AfterEach(func() {
			p = nil
		})
		It("failed the unmarshal", func() {
			templ := values.Templates{
				Values:   config.GetRaw(ci.PlaySpec{}),
				Internal: config.GetConfigRaw(),
			}
			valueBuf := new(bytes.Buffer)
			err := templ.GetValues(valueBuf)
			Expect(err).ToNot(Succeed())
		})
		It("get a good values.yaml", func() {
			templ := values.Templates{
				Values:   config.GetRaw(p.Spec),
				Internal: config.GetConfigRaw(),
			}
			valueBuf := new(bytes.Buffer)
			err := templ.GetValues(valueBuf)
			Expect(err).To(Succeed())
			Expect(valueBuf.String()).To(Equal(`---
env:
  - name: TEST
    value: "test"

lifecycle:
  enabled: true

image:
  repository: reg.example.com/group/repo
  tag: test

service:
  name: test-app

podLabels:
  application: test
dockerSecret:
  config: '{"auths":{"reg.example.com":{"auth":"dGVzdDp0ZXN0Cg=="}}}'

`))
		})
	})
})
