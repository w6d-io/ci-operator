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

import "bytes"

const (
	FileNameValues           string = "values.yaml"
	PostgresqlFileNameValues string = "values-postgresql.yaml"
	MongoDBFileNameValues    string = "values-mongodb.yaml"
)

// GetValues builds the values from the template from Play resource
func (in *Templates) GetValues(out *bytes.Buffer) error {

	if err := in.PrintTemplate(out, FileNameValues, HelmValuesTemplate); err != nil {
		return err
	}
	return nil
}

// GetMongoDBValues builds the values for mongoDB charts with dependencies elements from
// Play resource
func (in *Templates) GetMongoDBValues(out *bytes.Buffer) error {

	if err := in.PrintTemplate(out, MongoDBFileNameValues, MongoDBValuesTemplate); err != nil {
		return err
	}
	return nil
}

// GetPostgresValues builds the values for mongoDB charts with dependencies elements from
// Play resource
func (in *Templates) GetPostgresValues(out *bytes.Buffer) error {

	if err := in.PrintTemplate(out, PostgresqlFileNameValues, PostgresqlValuesTemplate); err != nil {
		return err
	}
	return nil
}
