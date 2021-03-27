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
	"bytes"
	"fmt"
	"github.com/go-logr/logr"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/minio"
	"github.com/w6d-io/ci-operator/internal/values"
)

func (wf *WFType) CreateValues(p *ci.Play, logger logr.Logger) error {
	log := logger.WithName("CreateValues")
	// get the task
	log.V(1).Info("create values.yaml")

	templ := values.Templates{
		Values:   config.GetRaw(p.Spec),
		Internal: config.GetConfigRaw(),
	}
	valueBuf := new(bytes.Buffer)
	if err := templ.GetValues(valueBuf); err != nil {
		return err
	}
	// TODO send Value to Minio or create Secret
	// put values.yaml in MinIO
	m := minio.New(logger)
	if m == nil {
		return fmt.Errorf("create minio instance return nil")
	}
	if err := m.PutString(logger, valueBuf.String(), BuildTarget(p, values.FileNameValues)); err != nil {
		return err
	}
	// TODO Create same process for MongoDB and PostgreSQL values
	// TODO for secret implementations update VolumeMount
	// TODO implements a method to factorized the process

	return nil
}

// BuildTarget return the path with project ID, pipeline ID and filename
func BuildTarget(play *ci.Play, filename string) (path string) {

	path = fmt.Sprintf("%v/%v/%s", play.Spec.ProjectID, play.Spec.PipelineID, filename)
	return
}
