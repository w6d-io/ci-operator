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
Created on 21/12/2020
*/

package values

import (
	"bytes"
	"context"
	"fmt"
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	"github.com/w6d-io/ci-operator/internal/util"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/speps/go-hashids"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/vault"
	"gopkg.in/yaml.v3"

	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	ValueLog = ctrl.Log.WithValues("package", "values")
	Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

func (in *Templates) PrintTemplate(ctx context.Context, out *bytes.Buffer, name string, templ string) error {
	log := ValueLog.WithName("PrintTemplate")
	correlationID, ok := util.GetCorrelationIDFromContext(ctx)
	if ok {
		log = log.WithValues("correlation_id", correlationID)
	}
	log.V(1).Info("templating")
	funcMap := sprig.TxtFuncMap()
	valuesMap := TxtFuncMap()

	t := template.Must(template.New(name).Funcs(funcMap).Funcs(valuesMap).Parse(templ))
	if err := t.Execute(out, in); err != nil {
		log.Error(err, "Execute")
		return err
	}
	return nil
}

func Vault(token string, path string, key string) (secret string) {
	log := ValueLog.WithName("Vault")
	vaultConfig := config.GetVault()
	if vaultConfig == nil || vaultConfig.GetHost() == "" {
		return
	}
	v := &vault.Config{
		Address: vaultConfig.GetHost(),
		Token:   token,
		Path:    path,
	}
	_ = v.GetSecret(ci.SecretKind(key), &secret, log)
	return
}

// HashID return hash from pid for the prefix url
func HashID(pid float64) (string, error) {
	hd := hashids.NewData()
	hd.Salt = config.GetHash().Salt
	hd.MinLength = config.GetHash().MinLength
	hd.Alphabet = Alphabet
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return "", fmt.Errorf("HashID failed with %s", err.Error())
	}
	var e string
	if e, err = h.Encode([]int{int(pid)}); err != nil {
		return "", fmt.Errorf("HashID failed with %s", err.Error())
	}
	return strings.ToLower(e), nil
}

// ToYaml returns the v interface into yaml format
func ToYaml(v interface{}) string {
	data, _ := yaml.Marshal(v)
	return string(data)
}

// TxtFuncMap returns a template FuncMap with HashID and ToYaml functions
func TxtFuncMap() template.FuncMap {
	return template.FuncMap{
		"hashID": HashID,
		"toYaml": ToYaml,
		"vault":  Vault,
	}
}
