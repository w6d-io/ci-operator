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
Created on 08/01/2021
*/

package minio

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/minio/minio-go/v6"
	"github.com/w6d-io/ci-operator/internal/config"
	"strings"
)

// New return a MinIO instance
func New(logger logr.Logger) Interface {
	log := logger.WithName("New").WithValues("package", "minio")
	log.V(1).Info(" instantiate minio client")

	if Instance != nil {
		return Instance
	}
	m := &Minio{
		Config: config.GetMinio(),
	}
	minioClient, err := minio.New(m.Config.Host, m.Config.AccessKey, m.Config.SecretKey, false)
	if err != nil {
		log.Error(err, "get minio client")
		return nil
	}
	if minioClient == nil {
		return nil
	}
	found, err := minioClient.BucketExists(m.Config.Bucket)
	if err != nil {
		log.Error(err, "check bucket exists")
		return nil
	}
	if !found {
		if err := minioClient.MakeBucket(m.Config.Bucket, "us-east-1"); err != nil {
			log.Error(err, "Making bucket", "bucket", m.Config.Bucket)
			return nil
		}
	}
	m.Client = minioClient
	log.V(1).Info("return minio client")
	Instance = m
	return Instance
}

// PutFile creates the object target in a bucket, with contents from file at source
func (m *Minio) PutFile(logger logr.Logger, source, target string) error {
	log := logger.WithName("PutFile").WithValues("package", "minio")

	log.V(1).Info("uploading", "source", source, "target", target)
	if _, err := m.Client.FPutObject(m.Config.Bucket, target, source, minio.PutObjectOptions{}); err != nil {
		log.Error(err, "FPutObject failed")
		return err
	}
	log.V(1).Info("upload done", "source", source, "target", target)
	return nil
}

// PutString creates the object target in a bucket, with the content of data of string type
func (m *Minio) PutString(logger logr.Logger, data, target string) error {
	log := logger.WithName("PutString").WithValues("package", "minio")
	r := strings.NewReader(data)
	log.V(1).Info("put object", "target", target)
	if _, err := m.Client.PutObject(m.Config.Bucket, target, r, r.Size(), minio.PutObjectOptions{}); err != nil {
		log.Error(err, "put data", "target", target)
		return fmt.Errorf("backup failed")
	}
	log.V(1).Info("upload done", "target", target)
	return nil
}
