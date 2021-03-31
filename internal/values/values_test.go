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
	"context"
	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/values"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

//docker_url: reg.example.com/group/repo:test

var playSpec = `
commit:
  message: no message
  ref: main
  sha: 3d5531613d5531613d5531613d5531613d553161
environment: staging
name: test
pipeline_id: 1
project_id: 1
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
				Client:   k8sClient,
				Values:   config.GetRaw(ci.PlaySpec{}),
				Internal: config.GetConfigRaw(),
			}
			valueBuf := new(bytes.Buffer)
			ctx := context.WithValue(context.Background(), "correlation_id", "unit-test")
			err := templ.GetValues(ctx, valueBuf, ctrl.Log)
			Expect(err).ToNot(Succeed())
		})
		It("get a good values.yaml", func() {
			templ := values.Templates{
				Client:   k8sClient,
				Values:   config.GetRaw(p.Spec),
				Internal: config.GetConfigRaw(),
			}
			valueBuf := new(bytes.Buffer)
			ctx := context.WithValue(context.Background(), "correlation_id", "unit-test")
			err := templ.GetValues(ctx, valueBuf, ctrl.Log)
			Expect(err).To(Succeed())
			Expect(valueBuf.String()).To(Equal(`---
env:
  - name: TEST
    value: "test"

serviceAccount: sa-1

lifecycle:
  enabled: true

image:
  repository: reg-ext.w6d.io/cxcm/1/test
  tag: 3d553161-main

service:
  name: test-app

podLabels:
  application: test
dockerSecret:
  config: '{"auths":{"reg.example.com":{"auth":"dGVzdDp0ZXN0Cg=="}}}'

`))
		})
		It("get a good values.yaml", func() {
			p.Spec.DockerURL = "reg.example.com/group/repo:test"
			templ := values.Templates{
				Client:   k8sClient,
				Values:   config.GetRaw(p.Spec),
				Internal: config.GetConfigRaw(),
			}
			valueBuf := new(bytes.Buffer)
			ctx := context.WithValue(context.Background(), "correlation_id", "unit-test")
			err := templ.GetValues(ctx, valueBuf, ctrl.Log)
			Expect(err).To(Succeed())
			Expect(valueBuf.String()).To(Equal(`---
env:
  - name: TEST
    value: "test"

serviceAccount: sa-1

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
		It("get values from configmap", func() {
			config.SetNamespace("default")
			ctx := context.WithValue(context.Background(), "correlation_id", "unit-test")
			cml := &corev1.NamespaceList{}
			err := k8sClient.List(ctx, cml)
			Expect(err).To(Succeed())
			Expect(len(cml.Items)).ToNot(Equal(0))
			var names []string
			for _, ns := range cml.Items {
				names = append(names, ns.Name)
			}
			//Expect(names).To(Equal([]string{}))
			cm := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "unit-test",
					Namespace: config.GetNamespace(),
				},
				Data: map[string]string{
					"values.yaml": "Test",
				},
			}
			err = k8sClient.Create(ctx, cm)
			Expect(err).To(Succeed())
			val := values.LookupOrDefaultValues(ctx, k8sClient, "deploy", values.HelmValuesTemplate)
			Expect(val).To(Equal("Test"))
			err = k8sClient.Delete(ctx, cm)
			Expect(err).To(Succeed())
		})
	})
})
