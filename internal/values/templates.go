/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 21/12/2020
*/

package values

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/ghodss/yaml"
	"github.com/speps/go-hashids"
)

// HelmValuesTemplate is the yaml template for helm command.
// render values.yaml according to the play resource contain
var HelmValuesTemplate = `---
{{- range $task := .Values.tasks }}
{{- range $key, $var := $task }}
{{- if and (eq $key "deploy") $var.variables }}
env:
{{- range $name, $value := $var.variables }}
  - name: {{ $name }}
    value: {{ $value | quote }}
{{- end }}
{{- end -}}
{{- end -}}
{{- end }}

service:
  name: {{ .Values.name }}-app

podLabels:
  application: {{ .Values.name }}

{{- if .Values.domain }}
ingress:
  enabled: true
  class: {{ .Internal.class }}
  host: {{ .Values.domain }}
{{- end }}

{{- if .Values.Dependencies }}
secrets:
{{- range $db := .Values.Dependencies }}
{{- range $name, $value := $db.Variables }}
  - name: {{ $name }}
    value: {{ $value | quote }}
    key: {{ $name | lower }}
	kind: env
{{- end }}
{{- end }}
{{- end }}

`

// MongoDBValuesTemplate MongoDB values chart template
var MongoDBValuesTemplate = `---
architecture: "standalone"
replicaCount: {{ default 1 .MongoDB.Replicas }}
auth:
  enabled: true
  rootPassword: {{ .MongoDB.RootPassword }}
  password:     {{ .MongoDB.Password}}
  username:     {{ .MongoDB.Username }}
  database:     {{ .MongoDB.Database }}
persistence:
  enabled: true
  size: {{ default 5Gi .MongoDB.Size }}
arbiter:
  enabled: false
metrics:
  enabled: true
  serviceMonitor:
    enabled: true
    namespace: monitoring
`

// PostgresqlValuesTemplate PostgreSQL values chart template
var PostgresqlValuesTemplate = `---
{{- $pass := randAlphaNum 20 }}
global:
  postgresql:
    postgresqlDatabase: {{ .Postgres.Database }}
    postgresqlUsername: {{ .Postgres.Username }}
    postgresqlPassword: {{ default $pass .Postgres.Password }}
postgresqlPostgresPassword: {{ .Postgres.PostgresPassword }}
persistence:
  enabled: true
  size: {{ default 5Gi .Postgres.Size }}
metrics:
  enabled: true
  serviceMonitor:
    enabled: true
`

var (
	Salt      string
	MinLength int
)

func (in *Templates) PrintTemplate(out *bytes.Buffer, name string, templ string) error {

	funcMap := sprig.TxtFuncMap()
	valuesMap := TxtFuncMap()

	t := template.Must(template.New(name).Funcs(funcMap).Funcs(valuesMap).Parse(templ))
	if err := t.Execute(out, in); err != nil {
		return err
	}
	return nil
}

// HashID return hash from pid for the prefix url
func HashID(pid float64) (string, error) {
	hd := hashids.NewData()
	hd.Salt = Salt
	hd.MinLength = MinLength
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
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return string(data)
}

// TxtFuncMap returns a template FuncMap with HashID and ToYaml functions
func TxtFuncMap() template.FuncMap {
	return template.FuncMap{
		"hashID": HashID,
		"toYaml": ToYaml,
	}
}
