/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 06/01/2021
*/

package secrets

import (
	"bytes"
	"context"
	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal/config"
	"github.com/w6d-io/ci-operator/internal/util"
	"github.com/w6d-io/ci-operator/internal/values"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

const (
	// filename use with s3cmd
	MinIOSecretKey string = ".s3cfg"
	// Prefix use for name of resource
	MinIOPrefixSecret string = "minio"
	// s3cfg config template
	MinIOSecretTemplate string = `
[default]
access_key = {{ .Values.access_key }}
secret_key = {{ .Values.secret_key }}
bucket_location = us-east-1
host_base = {{ .Values.host }}
host_bucket = {{ .Values.host }}/%(bucket)
default_mime_type = binary/octet-stream
enable_multipart = True
multipart_max_chunks = 10000
multipart_chunk_size_mb = 128
recursive = True
recv_chunk = 65536
send_chunk = 65536
server_side_encryption = False
signature_v2 = False
socket_timeout = 300
use_mime_magic = True
use_https = False
verbosity = WARNING
website_endpoint = {{ default "http" .Values.scheme }}://{{ .Values.host }}
`
)

// MinIOCreate creates the minio secret that contains the .s3cfg configuration
func (s *Secret) MinIOCreate(ctx context.Context, r client.Client, log logr.Logger) error {
	log = log.WithName("Create").WithValues("action", MinIOPrefixSecret)
	log.V(1).Info("creating")

	buf := new(bytes.Buffer)
	tpl := values.Templates{
		Values: config.GetMinioRaw(),
	}
	tpl.PrintTemplate(buf, MinIOSecretKey, MinIOSecretTemplate)
	namespacedName := util.GetCINamespacedName(MinIOPrefixSecret, s.Play)
	resource := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        namespacedName.Name,
			Namespace:   namespacedName.Namespace,
			Annotations: make(map[string]string),
			Labels:      util.GetCILabels(s.Play),
		},
		StringData: map[string]string{
			MinIOSecretKey: buf.String(),
		},
		Type: corev1.SecretTypeOpaque,
	}

	resource.Annotations[config.ScheduledTimeAnnotation] = time.Now().Format(time.RFC3339)
	if err := controllerutil.SetControllerReference(s.Play, resource, s.Scheme); err != nil {
		return err
	}
	if err := r.Create(ctx, resource); err != nil {
		return err
	}
	return nil
}
