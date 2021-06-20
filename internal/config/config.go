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
	"os"
	"sort"
	"strings"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/ghodss/yaml"
	"github.com/w6d-io/hook"
)

var (
	configLog = ctrl.Log.WithName("config")
)

const (
	DeployPrefixDefault = "cx"
	EnvPrefixDefault    = "W6D"
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
	for _, wh := range config.Webhooks {
		if wh.URLRaw != "" {
			if err := hook.Subscribe(wh.URLRaw, wh.Scope); err != nil {
				log.Error(err, "subscription failed")
				return err
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
		//{config.Volume.Name, "volume.name"},
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

// PodTemplate returns the tekton PodTemplate
func PodTemplate() *tkn.PodTemplate {
	if config.PodTemplate == nil {
		config.PodTemplate = new(tkn.PodTemplate)
	}
	return config.PodTemplate
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
	if config.DeployPrefix == "" {
		return DeployPrefixDefault
	}
	return config.DeployPrefix
}

// SetDeployPrefix record the prefix to use for deploy namespace name
func SetDeployPrefix(prefix string) {
	config.DeployPrefix = prefix
}

// GetMinio return the Minio structure
func GetMinio() *Minio {
	if config.Minio != nil {
		return config.Minio
	}
	return &Minio{}
}

// GetMinioRaw returns the Minio structure
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

func SetNamespace(namespace string) {
	config.Namespace = namespace
}

// GetWebhooks returns the list of url where to send the event
func GetWebhooks() []Webhook {
	return config.Webhooks
}

// GetValues return ValuesRef
func GetValues() Values {
	return config.ValuesRef
}

func GetHash() *Hash {
	if config.Hash != nil {
		return config.Hash
	}
	config.Hash = &Hash{}
	return config.Hash
}

// GetSalt return salt
func (h *Hash) GetSalt() string {
	return h.Salt
}

// GetMinLength return min hash length
func (h *Hash) GetMinLength() int {
	return h.MinLength
}

// GetVault return vault
func GetVault() *Vault {
	return config.Vault
}

// GetToken return the vault token
func (v *Vault) GetToken() string {
	return v.Token
}

// GetHost return the vault host
func (v *Vault) GetHost() string {
	return v.Host
}

// SetEnvPrefix record the prefix for environment variable
func SetEnvPrefix(prefix string) {
	config.EnvPrefix = prefix
}

// GetEnvPrefix return the prefix for variable environment
func GetEnvPrefix(elements ...string) string {
	prefix := EnvPrefixDefault
	if config.EnvPrefix != "" {
		prefix = config.EnvPrefix
	}
	toAdd := strings.Join(elements, "_")
	if toAdd != "" {
		//if toAdd[0] != '_' {
		//	toAdd = "_" + toAdd
		//}
		if toAdd[len(toAdd)-1] != '_' {
			toAdd += "_"
		}
	}
	return ToSnakeUpperCase(prefix + "_" + toAdd)
}

//var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
//var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeUpperCase(str string) string {
	//snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	//snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	snake := strings.ReplaceAll(str, "-", "_")
	return strings.ToUpper(snake)
}
