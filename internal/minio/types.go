/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 08/01/2021
*/

package minio

import (
	"github.com/minio/minio-go/v6"
	"github.com/w6d-io/ci-operator/internal/config"
)

//
type Minio struct {
	Client *minio.Client
	Config *config.Minio
}
