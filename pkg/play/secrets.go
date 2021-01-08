/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 24/12/2020
*/

package play

import (
	"github.com/go-logr/logr"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal"
	"github.com/w6d-io/ci-operator/internal/k8s/secrets"
	"github.com/w6d-io/ci-operator/internal/util"
)

// GitSecret fills the secret structure and add the git create method in the run list
func (wf *WFType) GitSecret(p *ci.Play, logger logr.Logger) error {
	log := logger.WithName("GitSecret").WithValues("cx-namespace",
		util.InNamespace(p))
	log.V(1).Info("Build git secret")
	secret := &secrets.Secret{
		WorkFlowStruct: internal.WorkFlowStruct{
			Scheme: wf.Scheme,
			Play:   p}}
	if err := wf.Add(secret.GitCreate); err != nil {
		return err
	}
	return nil
}

// DockerCredSecret fills the secret structure and add the Docker credential create method in the run list
func (wf *WFType) DockerCredSecret(p *ci.Play, logger logr.Logger) error {
	log := logger.WithName("DockerCredSecret").WithValues("cx-namespace",
		util.InNamespace(p))
	log.V(1).Info("Build docker credential secret")
	secret := &secrets.Secret{
		WorkFlowStruct: internal.WorkFlowStruct{
			Scheme: wf.Scheme,
			Play:   p}}
	if err := wf.Add(secret.DockerCredCreate); err != nil {
		return err
	}
	return nil
}

// MinIOSecret fills the secret structure to add minio create method in the run list
func (wf *WFType) MinIOSecret(p *ci.Play, logger logr.Logger) error {
	log := logger.WithName("MinIOSecret").WithValues("cx-namespace",
		util.InNamespace(p))
	log.V(1).Info("Build minio secret")
	secret := &secrets.Secret{
		WorkFlowStruct: internal.WorkFlowStruct{
			Scheme: wf.Scheme,
			Play:   p}}
	if err := wf.Add(secret.MinIOCreate); err != nil {
		return err
	}
	return nil
}
