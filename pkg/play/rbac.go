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
Created on 07/01/2021
*/

package play

import (
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"

	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal"
	"github.com/w6d-io/ci-operator/internal/k8s/rbac"
)

func (wf *WFType) Rbac(p *ci.Play, logger logr.Logger) error {
	log := logger.WithName("Rbac")
	log.V(1).Info("Build CI role-binding")

	resourceCI := &rbac.CI{
		WorkFlowStruct: internal.WorkFlowStruct{
			Play:   p,
			Scheme: wf.Scheme,
		},
	}

	if err := wf.Add(resourceCI.Create); err != nil {
		log.V(1).Info("Build Deploy role-binding")
		return err
	}
	resourceDeploy := &rbac.Deploy{
		WorkFlowStruct: internal.WorkFlowStruct{
			Play:   p,
			Scheme: wf.Scheme,
		},
	}

	if err := wf.Add(resourceDeploy.Create); err != nil {
		return err
	}
	return nil
}
