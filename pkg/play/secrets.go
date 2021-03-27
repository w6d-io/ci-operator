/*
Copyright 2020 WILDCARD

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
Created on 24/12/2020
*/

package play

import (
	"github.com/go-logr/logr"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal"
	"github.com/w6d-io/ci-operator/internal/k8s/secrets"
)

// GitSecret fills the secret structure and add the git create method in the run list
func (wf *WFType) GitSecret(p *ci.Play, logger logr.Logger) error {
	log := logger.WithName("GitSecret")
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
	log := logger.WithName("DockerCredSecret")
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
	log := logger.WithName("MinIOSecret")
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
