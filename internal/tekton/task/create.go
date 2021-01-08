/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 07/01/2021
*/

package task

import (
	"context"

	ci "github.com/w6d-io/ci-operator/api/v1alpha1"

	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (t *Task) Parse(ctx context.Context, logger logr.Logger) error {
	log := logger.WithName("Parse").WithValues("cx-namespace", util.InNamespace(t.Play))

	for pos, m := range t.Play.Spec.Tasks {
		for name := range m {
			switch name {
			case ci.Build:
				// TODO call method build
				log.WithValues("task", ci.Build).V(1).Info("launch")
				t.Index = pos
				if err := t.Build(ctx, logger); err != nil {
					return err
				}
			case ci.Sonar:
				// TODO call methods for unitTest
				log.WithValues("task", ci.Sonar).V(1).Info("launch")
				t.Index = pos
				if err := t.Sonar(ctx, logger); err != nil {
					return err
				}
			case ci.UnitTests:
				// TODO call methods for unitTest
				log.WithValues("task", ci.UnitTests).V(1).Info("launch")
				t.Index = pos
				if err := t.UnitTest(ctx, logger); err != nil {
					return err
				}
			case ci.IntegrationTests:
				// TODO call methods for integrationTest
				log.WithValues("task", ci.IntegrationTests).V(1).Info("launch")
				t.Index = pos
				if err := t.IntTest(ctx, logger); err != nil {
					return err
				}
			case ci.Deploy:
				// TODO call methods for deploy
				log.WithValues("task", ci.Deploy).V(1).Info("launch")
				t.Index = pos
				if err := t.Deploy(ctx, logger); err != nil {
					return err
				}
			case ci.Clean:
				// TODO call methods for clean
				log.WithValues("task", ci.Clean).V(1).Info("launch")
				t.Index = pos
				if err := t.Clean(ctx, logger); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (t *Task) Add(taskFunc func(context.Context, client.Client, logr.Logger) error) error {
	t.Creates = append(t.Creates, taskFunc)
	return nil
}
