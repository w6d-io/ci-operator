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
	"errors"
	"github.com/go-logr/logr"
	"github.com/hashicorp/vault/api"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
)

// GetSecret returns the secret data from vault
func (c *Config) GetSecret(key ci.SecretKind, rec interface{}, log logr.Logger) error {
	log = log.WithName("GetSecret")
	log.V(1).Info("get vault secret")

	client, err := api.NewClient(&api.Config{Address: c.Address, HttpClient: httpClient})
	if err != nil {
		return err
	}
	client.SetToken(c.Token)
	log.V(1).Info("read data", "path", c.Path)
	data, err := client.Logical().Read(c.Path)
	if err != nil {
		log.Error(err, "read data from vault")
		return err
	}
	log.V(1).Info("data", "path", c.Path, "data", data)
	if data == nil {
		log.Error(nil, "data from vault is empty")
		return errors.New("data from vault is empty")
	}
	sec, ok := data.Data[string(key)]
	if !ok {
		log.Error(nil, "data from vault not contains the key")
		return errors.New("data from vault not contains the key")
	}
	b, err := json.Marshal(sec)
	if err != nil {
		log.Error(err, "marshal data from vault")
		return err
	}

	if err := json.Unmarshal(b, rec); err != nil {
		return err
	}
	return nil
}
