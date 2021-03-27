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
Created on 08/12/2020
*/

package play

import (
	"github.com/go-logr/logr"

	ci "github.com/w6d-io/ci-operator/api/v1alpha1"

	intpipeline "github.com/w6d-io/ci-operator/internal/tekton/pipeline"
)

func (wf *WFType) SetPipeline(p *ci.Play, logger logr.Logger) error {
	log := logger.WithName("SetPipeline")
	log.Info("Build pipeline")
	pipeline := &intpipeline.Pipeline{
		Play:   p,
		Scheme: wf.Scheme,
	}
	if err := pipeline.Parse(log); err != nil {
		return err
	}

	log.V(1).Info("add pipeline create method")
	if err := wf.Add(pipeline.Create); err != nil {
		return err
	}
	return nil
}
