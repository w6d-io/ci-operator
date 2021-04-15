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
Created on 26/02/2021
*/
package http

import (
    "bytes"
    "encoding/json"
    "github.com/avast/retry-go"
    "io/ioutil"
    "net/http"
    "net/url"
    "strconv"
    "time"
)

func (h *HTTP) Send(payload interface{}, URL *url.URL) error {
    log := logger.WithName("Send").WithValues("URL", URL.Redacted())
    h.Username    = URL.User.Username()
    h.Password, _ = URL.User.Password()
    query := URL.Query()
    to, ok := query["timeout"]
    client := http.Client{
        Timeout: 5 * time.Second,
    }
    if ok {
        n, err := strconv.ParseInt(to[0], 10, 64)
        if err != nil {
            log.Error(err, "convert timeout failed")
            return err
        }
        client = http.Client{
            Timeout: time.Duration(n) * time.Second,
        }
    }
    log.V(1).Info("marshal payload")
    data, err := json.Marshal(payload)
    if err != nil {
        log.Error(err, "marshal failed")
        return err
    }
    if err := retry.Do(
        func() error {
            log.V(1).Info("post payload", "attempt", retry.DefaultAttempts)
            response, err := client.Post(URL.String(), "application/json", bytes.NewBuffer(data))
            if err == nil {
                defer func() {
                    if err := response.Body.Close(); err != nil {
                        log.Error(err, "close http response")
                        return
                    }
                }()
                body, err := ioutil.ReadAll(response.Body)
                if err != nil {
                    log.Error(err, "get response body")
                }
                log.Info(string(body))
                return nil
            }
            log.Error(err, "post data failed")
            return err
        },
        retry.Attempts(5),
        ); err != nil {
        return err
    }

    return nil
}

func (HTTP) Validate(_ *url.URL) error {
    return nil
}