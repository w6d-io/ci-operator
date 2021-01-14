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
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/ghodss/yaml"
	"github.com/speps/go-hashids"
)

var (
	Salt      string
	MinLength int
)

var (
	ValueLog = ctrl.Log.WithValues("package", "values")
)

func (in *Templates) PrintTemplate(out *bytes.Buffer, name string, templ string) error {
	log := ValueLog.WithName("PrintTemplate")
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
