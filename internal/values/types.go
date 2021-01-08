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
	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
)

// Templates structure with methods for templating
// - values.yaml
// - values.mongodb.yaml
// - values.postgre.yaml
// - .s3cfg
type Templates struct {
	Play     *ci.Play
	Spec     ci.PlaySpec
	Internal map[string]interface{}
	Values   map[string]interface{}
}
