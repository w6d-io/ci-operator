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
	"time"

	ci "github.com/w6d-io/ci-operator/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/avast/retry-go"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
)

var (
	logger = ctrl.Log.WithName("WebHook")
)

type PlayPayload struct {
	Object     Object    `json:"object,omitempty"`
	ProjectID  int64     `json:"project_id,omitempty"`
	PipelineID int64     `json:"pipeline_id,omitempty"`
	RepoURL    string    `json:"repo_url,omitempty"`
	Commit     ci.Commit `json:"ref,omitempty"`
	Stack      ci.Stack  `json:"stack,omitempty"`
	Status     ci.State  `json:"status,omitempty"`
}

type Object struct {
	NamespacedName types.NamespacedName `json:"namespaced_name,omitempty"`
	Kind           string               `json:"kind,omitempty"`
}

func BuildPlayPayload(play *ci.Play, status ci.State, logger logr.Logger) error {
	log := logger.WithName("BuildPayload")
	log.V(1).Info("build payload")

	payload = &PlayPayload{
		Object: Object{
			Kind: "pipelinerun",
		},
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

// GetStatus returns the play status from payload
func (p *PlayPayload) GetStatus() ci.State {
	return p.Status
}

// SetStatus sets the play status to payload
func (p *PlayPayload) SetStatus(state ci.State) {
	p.Status = state
}

// GetObjectName returns the object name from payload
func (p *PlayPayload) GetObjectNamespacedName() types.NamespacedName {
	return p.Object.NamespacedName
}

// SetObjectName sets the object name in payload
func (p *PlayPayload) SetObjectNamespacedName(namespacedName types.NamespacedName) {
	p.Object.NamespacedName = namespacedName
}

func (p *PlayPayload) Send(URL string) error {
	if URL == "" {
		return nil
	}
	log := logger.WithName("Send").WithValues("URL", URL)
	log.V(1).Info("create http client")
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	log.V(1).Info("marshal payload")
	data, err := json.Marshal(p)
	if err != nil {
		log.Error(err, "json conversion failed")
		return err
	}
	if err := retry.Do(
		func() error {
			log.V(1).WithValues("send", retry.DefaultAttempts).Info("post payload")
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
		retry.Attempts(5),
	); err != nil {
		return err
	}
	return nil
}

func (p *PlayPayload) DoSend(whs []Webhook) error {
	log := logger.WithName("DoSend")
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
	for range whs {
		if err := <-errc; err != nil {
			log.Error(err, "Send failed")
			return err
		}
	}
	return nil
}
