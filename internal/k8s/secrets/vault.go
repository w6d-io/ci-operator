/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 16/04/2021
*/

package secrets

import (
	"github.com/go-logr/logr"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/vault"
)

// VaultCreate creating all secrets from vault according the paths and the
// token scope pass through play resource
//func (s *Secret) VaultCreate(ctx context.Context, r client.Client, logger logr.Logger) error {
//	log := logger.WithName("Create").WithValues("action", "VaultSecret")
//	log.V(1).Info("creating")
//	var key ci.SecretKind
//	var vaultSecret ci.VaultSecret
//	for key, vaultSecret = range s.Play.Spec.Vault.Secrets {
//		data := s.GetVaultSecret(string(key), vaultSecret, log)
//		namespacedName := util.GetCINamespacedName(string(key), s.Play)
//		resource := &corev1.Secret{
//			ObjectMeta: metav1.ObjectMeta{
//				Name:        namespacedName.Name,
//				Namespace:   namespacedName.Namespace,
//				Labels:      util.GetCILabels(s.Play),
//				Annotations: make(map[string]string),
//			},
//			StringData: map[string]string{
//				string(key): data,
//			},
//			Type: corev1.SecretTypeOpaque,
//		}
//		resource.Annotations[config.ScheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
//		if err := controllerutil.SetControllerReference(s.Play, resource, s.Scheme); err != nil {
//			return err
//		}
//		if err := r.Create(ctx, resource); err != nil {
//			log.Error(err, "create")
//			return err
//		}
//		if err := sa.Update(ctx, resource.Name,
//			util.GetCINamespacedName(sa.Prefix, s.Play), r); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}

// GetVaultSecret return the secret from vault
func (s *Secret) GetVaultSecret(key string, sec ci.VaultSecret, logger logr.Logger) (secret string) {
	log := logger.WithName("GetVaultSecret").WithValues("path", sec.Path)
	log.V(1).Info("look at")
	vaultConfig := config.GetVault()
	if vaultConfig == nil || vaultConfig.GetHost() == "" {
		log.Error(nil, "vault is not configured")
		return
	}
	if sec.Path == "" {
		log.Error(nil, "vault path is empty")
		return
	}
	log.V(1).Info("get secret from vault")
	v := &vault.Config{
		Address: vaultConfig.GetHost(),
		Token:   s.GetToken(vaultConfig, logger),
		Path:    sec.Path,
	}
	if err := v.GetSecret(key, &secret, logger); err != nil {
		log.Error(err, "get secret from vault failed")
		return
	}
	logger.V(1).Info("return vault secret")
	return
}

// GetToken from play is exist or get it from config
func (s *Secret) GetToken(v *config.Vault, logger logr.Logger) string {
	logger.V(1).Info("get token")
	if s.Play.Spec.Vault != nil && s.Play.Spec.Vault.Token != "" {
		logger.V(1).Info("return play token")
		return s.Play.Spec.Vault.Token
	}
	if v != nil {
		logger.V(1).Info("return config vault token")
		return v.Token
	}
	return ""
}
