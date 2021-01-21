/*
Copyright 2020 WILDCARD SA.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
Created on 20/01/2021
*/
package webhook

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/avast/retry-go"
	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-operator/internal/config"
	"time"
)

var (
	logger = ctrl.Log.WithName("WebHook")
)

func BuildPlayPayload(play *ci.Play, status ci.State, logger logr.Logger) error {
	log := logger.WithName("BuildPayload")
	log.V(1).Info("build payload")

	payload = &PlayPayload{
		ObjectKind: play.APIVersion,
		ProjectID:  play.Spec.ProjectID,
		PipelineID: play.Spec.PipelineID,
		RepoURL:    play.Spec.RepoURL,
		Commit:     play.Spec.Commit,
		Stack:      play.Spec.Stack,
		Status:     status,
	}
	return nil
}

// GetPayLoad returns the
func GetPayLoad() Payload {
	return payload
}

func (p *PlayPayload) SetStatus(state ci.State) {
	p.Status = state
}

func (p *PlayPayload) Send(URL string) error {
	log := logger.WithName("Send")
	if URL == "" {
		return nil
	}
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	data, err := json.Marshal(p)
	if err != nil {
		log.Error(err, "json conversion failed")
		return err
	}
	retry.DefaultAttempts = 5
	if err := retry.Do(
		func() error {
			response, err := client.Post(URL, "application/json", bytes.NewBuffer(data))
			if err == nil {
				defer func() {
					if err := response.Body.Close(); err != nil {
						log.Error(err, "close http response")
						return
					}
				}()
				data, err := ioutil.ReadAll(response.Body)
				if err != nil {
					log.Error(err, "get response body")
				}
				log.Info(string(data))
			}
			log.Error(err, "Post data returns failed")
			return err
		},
	); err != nil {
		return err
	}
	return nil
}

func (p *PlayPayload) DoSend() error {
	log := logger.WithName("DoSend")
	whs := config.GetWebhooks()
	errc := make(chan error, len(whs))
	quit := make(chan struct{})
	defer close(quit)

	for _, wh := range whs {
		go func(URL string) {
			logg := log.WithValues("url", URL)
			select {
			case errc <- p.Send(URL):
				logg.Info("sent")
			case <-quit:
				logg.Info("quit")
			}
		}(wh.URLRaw)
	}
	for range config.GetWebhooks() {
		if err := <-errc; err != nil {
			log.Error(err, "Send failed")
			return err
		}
	}
	return nil
}
