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
Created on 30/12/2020
*/

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"sort"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/ghodss/yaml"
	"github.com/w6d-io/ci-operator/internal/values"
)

var (
	configLog = ctrl.Log.WithName("config")
)

// New get the filename and fill Config struct
func New(filename string) error {
	// TODO add dynamic configuration feature
	log := ctrl.Log.WithName("controllers").WithName("Config")
	log.V(1).Info("read config file")
	config = new(Config)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error(err, "error reading the configuration")
		return err
	}
	if err := yaml.Unmarshal(data, config); err != nil {
		log.Error(err, "Error unmarshal the configuration")
		return err
	}
	config.Namespace = os.Getenv("NAMESPACE")
	if config.Volume.Name == "" {
		config.Volume = tkn.WorkspaceBinding{
			Name:     "ws",
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		}
	}
	if config.DeployPrefix == "" {
		config.DeployPrefix = "cx"
	}
	values.Salt = config.Hash.Salt
	values.MinLength = config.Hash.MinLength
	for i, wh := range config.Webhooks {
		if wh.URLRaw != "" {
			config.Webhooks[i].URL, err = url.Parse(wh.URLRaw)
			if err != nil {
				return fmt.Errorf("webhook (%s) config error %s", wh.Name, err)
			}
		}
	}
	return nil
}

// Validate returns an error if is mandatory Config missing
func Validate() error {
	var unset []string
	for _, f := range []struct {
		v, name string
	}{
		{config.DefaultDomain, "domain"},
		{config.Ingress.Class, "ingress.class"},
		{config.Volume.Name, "volume.name"},
	} {
		if f.v == "" {
			unset = append(unset, f.name)
		}
	}
	if len(unset) > 0 {
		sort.Strings(unset)
		return fmt.Errorf("found unset config flags: %s", unset)
	}
	return nil
}

// Workspaces returns list of Tekton WorkspaceDeclaration
func Workspaces() []tkn.WorkspaceDeclaration {
	return config.Workspaces
}

// Volume returns Tekton WorkspaceBinding
func Volume() tkn.WorkspaceBinding {
	return config.Volume
}

// GetWorkspacePath returns path from workspace
func GetWorkspacePath(name string) string {
	for _, wk := range config.Workspaces {
		if wk.Name == name {
			subPath := "/workspaces/" + wk.Name
			if wk.MountPath != "" {
				subPath = wk.MountPath
			}
			return subPath
		}
	}
	return ""
}

// GetConfig returns the Config values
func GetConfig() *Config {
	return config
}

// GetConfigRaw returns the Config structure in map[string]interface
func GetConfigRaw() map[string]interface{} {
	return GetRaw(config)
}

// GetClusterRole returns the ClusterRole from config
func GetClusterRole() string {
	return config.ClusterRole
}

// GetDeployPrefix returns the prefix to use for deploy namespace name
func GetDeployPrefix() string {
	return config.DeployPrefix
}

// GetMinio return the Minio structure
func GetMinio() *Minio {
	if config.Minio != nil {
		return config.Minio
	}
	return &Minio{}
}

// GetMinio returns the Minio structure
func GetMinioRaw() map[string]interface{} {
	if config.Minio != nil {
		return GetRaw(config.Minio)
	}
	return nil
}

// GetRaw return a map[string]interface from the interface
func GetRaw(input interface{}) map[string]interface{} {
	output := map[string]interface{}{}
	data, err := json.Marshal(input)
	if err != nil {
		configLog.Error(err, "GetRaw")
		return nil
	}
	if err := json.Unmarshal(data, &output); err != nil {
		return nil
	}
	return output
}

// GetHost method returns the host from Minio structure
func (m *Minio) GetHost() string {
	return m.Host
}

// GetAccessKey method returns the access_key from Minio structure
func (m *Minio) GetAccessKey() string {
	return m.AccessKey
}

// GetSecretKey method returns the secret_key from Minio structure
func (m *Minio) GetSecretKey() string {
	return m.AccessKey
}

// GetBucket method returns the minio bucket
func (m *Minio) GetBucket() string {
	if config.Minio != nil {
		return config.Minio.Bucket
	}
	return ""
}

func GetNamespace() string {
	return config.Namespace
}

// GetWebhooks returns the list of url where to send the event
func GetWebhooks() []Webhook {
	return config.Webhooks
}
