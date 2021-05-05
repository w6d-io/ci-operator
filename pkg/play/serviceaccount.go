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
Created on 28/12/2020
*/

package play

import (
	"github.com/go-logr/logr"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal"
	"github.com/w6d-io/ci-operator/internal/k8s/sa"
)

func (wf *WFType) ServiceAccount(p *ci.Play, logger logr.Logger) error {
	log := logger.WithName("GitSecret")
	log.V(1).Info("Build service account for CI")
	cisa := &sa.CI{
		WorkFlowStruct: internal.WorkFlowStruct{
			Play:   p,
			Scheme: wf.Scheme,
		},
	}
	if err := wf.Add(cisa.Create); err != nil {
		log.V(1).Info("Build service account for ci")
		return err
	}

	// Comment due to this serviceaccount should be created by helm during the deployment
	deploy := &sa.Deploy{
		WorkFlowStruct: internal.WorkFlowStruct{
			Play:   p,
			Scheme: wf.Scheme,
		},
	}
	log.V(1).Info("Build service account for Deploy")
	if err := wf.Add(deploy.Create); err != nil {
		log.V(1).Info("Build service account for deploy")
		return err
	}
	return nil
}
