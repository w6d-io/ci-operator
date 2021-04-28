/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 21/04/2021
*/

package secrets_test

import (
	"github.com/w6d-io/ci-operator/internal"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/k8s/secrets"
	corev1 "k8s.io/api/core/v1"

	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vault secret", func() {
	Context("", func() {
		It("Gets token", func() {
			s := &secrets.Secret{
				WorkFlowStruct: internal.WorkFlowStruct{
					Play: &ci.Play{
						Spec: ci.PlaySpec{
							Vault: &ci.Vault{},
						},
					},
				},
			}
			By("gets empty string on vault is nil")
			Expect(s.GetToken(nil, ctrl.Log)).To(Equal(""))

			v := &config.Vault{}
			By("gets empty string on vault token empty")
			Expect(s.GetToken(v, ctrl.Log)).To(Equal(""))

			v.Token = "token"
			By("gets vault token")
			Expect(s.GetToken(v, ctrl.Log)).To(Equal("token"))

			By("gets the play token")
			s.Play.Spec.Vault.Token = "play token"
			Expect(s.GetToken(v, ctrl.Log)).To(Equal("play token"))

		})
		It("Get Vault Secret", func() {
			s := &secrets.Secret{
				WorkFlowStruct: internal.WorkFlowStruct{
					Play: &ci.Play{
						Spec: ci.PlaySpec{
							Vault: &ci.Vault{
								Token: "play token",
							},
						},
					},
				},
			}
			By("gets empty string by empty config")
			sec := ci.VaultSecret{}
			Expect(s.GetVaultSecret(corev1.DockerConfigJsonKey, sec, ctrl.Log)).To(Equal(""))

			By("set vault config")
			Expect(config.New("testdata/file1.yaml")).To(Succeed())

			By("gets empty string by empty path")
			Expect(s.GetVaultSecret(corev1.DockerConfigJsonKey, sec, ctrl.Log)).To(Equal(""))

			By("gets empty string due to vault server absent")
			sec.Path = "test/test"
			Expect(s.GetVaultSecret(corev1.DockerConfigJsonKey, sec, ctrl.Log)).To(Equal(""))
		})
		It("Get Secret", func() {
			s := &secrets.Secret{
				WorkFlowStruct: internal.WorkFlowStruct{
					Play: &ci.Play{
						Spec: ci.PlaySpec{
							Vault: nil,
							Secret: map[string]string{
								corev1.DockerConfigJsonKey: "{}",
							},
						},
					},
				},
			}
			By("gets the play secret")
			Expect(s.GetSecret(corev1.DockerConfigJsonKey, ctrl.Log)).To(Equal("{}"))

			By("gets the play secret because vault does not have the key")
			s.Play.Spec.Vault = &ci.Vault{
				Token: "play token",
				Secrets: map[ci.SecretKind]ci.VaultSecret{
					"key": {Path: "test/test"}}}
			Expect(s.GetSecret(corev1.DockerConfigJsonKey, ctrl.Log)).To(Equal("{}"))

			s.Play.Spec.Vault = &ci.Vault{
				Token: "play token",
				Secrets: map[ci.SecretKind]ci.VaultSecret{
					corev1.DockerConfigJsonKey: {Path: "test/test"}}}
			By("gets the play secret due to vault failed")
			Expect(s.GetSecret(corev1.DockerConfigJsonKey, ctrl.Log)).To(Equal("{}"))
		})
	})
})
