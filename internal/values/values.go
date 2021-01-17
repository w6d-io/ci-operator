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
)

const (
	FileNameValues           string = "values.yaml"
	PostgresqlFileNameValues string = "values-postgresql.yaml"
	MongoDBFileNameValues    string = "values-mongodb.yaml"
)

// TODO replace or add these templates by configmap or a new resource

// HelmValuesTemplate is the yaml template for helm command.
// render values.yaml according to the play resource contain
var HelmValuesTemplate = `---
{{- $defaultDomain := printf "%v.%s" (.Values.project_id | hashID) .Internal.domain }}
{{- $repository := (printf "reg-ext.w6d.io/cxcm/%v/%v" .Values.project_id .Values.name) }}
{{- $tag := (printf "%v-%v" (substr 0 7 .Values.commit.sha) .Values.commit.ref) }}
{{- if .Values.docker_url }}
{{- $part := split ":" .Values.docker_url }}
{{- $repository = $part._0 }}
{{- $tag = $part._1 }}
{{- end }}
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

image:
  repository: {{ $repository }}
  tag: {{ $tag }}

service:
  name: {{ .Values.name }}-app

podLabels:
  application: {{ .Values.name }}

{{- if .Values.domain }}
ingress:
  enabled: true
  class: {{ .Internal.ingress.class }}
  host: {{ default $defaultDomain .Values.domain }}
  issuer: {{ .Internal.ingress.issuer | quote }}
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

// GetValues builds the values from the template from Play resource
func (in *Templates) GetValues(out *bytes.Buffer) error {
	log := ValueLog.WithName("GetValues")
	log.V(1).Info("templating")
	if err := in.PrintTemplate(out, FileNameValues, HelmValuesTemplate); err != nil {
		return err
	}
	return nil
}

// GetMongoDBValues builds the values for mongoDB charts with dependencies elements from
// Play resource
func (in *Templates) GetMongoDBValues(out *bytes.Buffer) error {
	log := ValueLog.WithName("GetMongoDBValues")
	log.V(1).Info("templating")
	if err := in.PrintTemplate(out, MongoDBFileNameValues, MongoDBValuesTemplate); err != nil {
		return err
	}
	return nil
}

// GetPostgresValues builds the values for mongoDB charts with dependencies elements from
// Play resource
func (in *Templates) GetPostgresValues(out *bytes.Buffer) error {
	log := ValueLog.WithName("GetPostgresValues")
	log.V(1).Info("templating")

	if err := in.PrintTemplate(out, PostgresqlFileNameValues, PostgresqlValuesTemplate); err != nil {
		return err
	}
	return nil
}
