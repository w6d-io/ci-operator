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
Created on 16/02/2021
*/
package kafka

import (
    "k8s.io/klog/klogr"
    "time"
)

var (
    logger = klogr.New()
)

type Kafka struct {
    Username        string
    Password        string
    BootstrapServer string
    Topic           string
}

// Option ...
type Option func(*Options)

// Options ...
type Options struct {
    Protocol          string
    Mechanisms        string
    Async             bool
    SessionTimeout    time.Duration
    MaxPollInterval   time.Duration
    WriteTimeout      time.Duration
    ReadTimeout       time.Duration
    BatchTimeout      time.Duration
    MaxWait           time.Duration
    StatInterval      time.Duration
    NumPartitions     int
    ReplicationFactor int
    MinBytes          int
    MaxBytes          int
    AuthKafka         bool
    FullStats         bool
    Debugs            []string
    GroupInstanceID   string
    ConfigMapKey      string
}
