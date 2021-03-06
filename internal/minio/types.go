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
	"github.com/go-logr/logr"
	"github.com/minio/minio-go/v6"
	"github.com/w6d-io/ci-operator/internal/config"
)

type Interface interface {
	PutFile(logr.Logger, string, string) error
	PutString(logr.Logger, string, string) error
}

// Minio contains the instance and the configuration
type Minio struct {
	Client *minio.Client
	Config *config.Minio
}

var Instance Interface
