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

package vault

import (
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/hashicorp/vault/api"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
)

// TODO add the toggle for vault enabling

// GetSecret returns the secret data from vault
func (c *Config) GetSecret(play ci.Play, rec interface{}, log logr.Logger) error {
	log = log.WithName("GitSecret")
	log.V(1).Info("get vault secret")

	client, err := api.NewClient(&api.Config{Address: c.Address, HttpClient: httpClient})
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%v/%v", Engine, play.Spec.ProjectID, play.Spec.PipelineID)
	data, err := client.Logical().Read(path)
	if err != nil {
		log.Error(err, "read data from vault")
		return err
	}
	b, err := json.Marshal(data.Data)
	if err != nil {
		log.Error(err, "marshal data from vault")
		return err
	}

	if err := json.Unmarshal(b, rec); err != nil {
		return err
	}
	return nil
}

func (c *Config) GetToken() {
	// TODO in progress
}

func (c *Config) SaveObject() {
	// TODO in progress

}
