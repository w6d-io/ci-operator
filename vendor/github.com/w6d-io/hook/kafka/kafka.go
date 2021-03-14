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
Created on 08/02/2021
*/
package kafka

import (
    "errors"
    "net/url"
    "time"

    "github.com/avast/retry-go"
)

func (k *Kafka) Send(payload interface{}, URL *url.URL) error {
    passwd, ok := URL.User.Password()
    query := URL.Query()
    k.BootstrapServer = URL.Host
    k.Topic           = query["topic"][0]
    k.Username        = URL.User.Username()
    k.Password        = passwd

    var (
        async      bool
        messageKey string
        protocol   = "SASL_SSL"
        mechanisms = "PLAIN"
    )
    async = len(query["async"]) > 0 && query["async"][0] == "true"
    if len(query["messagekey"]) > 0 {
        messageKey = query["messagekey"][0]
    }
    if len(query["protocol"]) > 0 {
        protocol = query["protocol"][0]
    }
    if len(query["mechanisms"]) > 0 {
        mechanisms = query["mechanisms"][0]
    }
    if err := retry.Do(
        func() error {
            if err := k.Producer(messageKey, payload,
                AuthKafka(ok), Async(async),
                Protocol(protocol), Mechanisms(mechanisms),
            ); err != nil {
                return err
            }
            return nil
        },
        retry.Attempts(5),
    ); err != nil {
        return err
    }
    //logger.V(1).Info("send payload by kafka", "payload", payload,
    //	"address", URL.Host)
    return nil
}

func (k *Kafka) Validate(URL *url.URL) error {
    if URL == nil {
        return nil
    }
    values := URL.Query()
    if _, ok := values["topic"]; !ok {
        logger.Error(errors.New("missing topic"), URL.Redacted())
        return errors.New("missing topic")
    }
    return nil
}

// NewOptions ...
func NewOptions(opts ...Option) Options {
    opt := Options{
        Protocol:          "SASL_SSL",
        Mechanisms:        "PLAIN",
        Async:             true,
        SessionTimeout:    10 * time.Second,
        MaxPollInterval:   5 * time.Minute,
        WriteTimeout:      10 * time.Second,
        ReadTimeout:       10 * time.Second,
        BatchTimeout:      1 * time.Millisecond,
        MaxWait:           2 * time.Millisecond,
        StatInterval:      5 * time.Second,
        MinBytes:          10e3,
        MaxBytes:          10e6,
        NumPartitions:     1,
        ReplicationFactor: 3,
        AuthKafka:         false,
        FullStats:         false,
        GroupInstanceID:   "",
        Debugs:            []string{},
        ConfigMapKey:      "kafka",
    }
    for _, o := range opts {
        o(&opt)
    }
    return opt
}

// Protocol option
func Protocol(p string) Option {
    return func(o *Options) {
        o.Protocol = p
    }
}

// Mechanisms option
func Mechanisms(m string) Option {
    return func(o *Options) {
        o.Mechanisms = m
    }
}

// Async option
func Async(b bool) Option {
    return func(o *Options) {
        o.Async = b
    }
}

// WriteTimeout option
func WriteTimeout(t time.Duration) Option {
    return func(o *Options) {
        o.WriteTimeout = t
    }
}

// MaxWait option
func MaxWait(t time.Duration) Option {
    return func(o *Options) {
        o.MaxWait = t
    }
}

// StatInterval option
func StatInterval(t time.Duration) Option {
    return func(o *Options) {
        o.StatInterval = t
    }
}

// MaxBytes option
//func MaxBytes(b int) Option {
//	return func(o *Options) {
//		o.MaxBytes = b
//	}
//}

// MinBytes option
//func MinBytes(b int) Option {
//	return func(o *Options) {
//		o.MinBytes = b
//	}
//}

// NumPartitions option
func NumPartitions(n int) Option {
    return func(o *Options) {
        o.NumPartitions = n
    }
}

// ReplicationFactor option
func ReplicationFactor(r int) Option {
    return func(o *Options) {
        o.ReplicationFactor = r
    }
}

// AuthKafka option
func AuthKafka(b bool) Option {
    return func(o *Options) {
        o.AuthKafka = b
    }
}

// FullStats option
func FullStats(b bool) Option {
    return func(o *Options) {
        o.FullStats = b
    }
}

// Debugs option
func Debugs(d []string) Option {
    return func(o *Options) {
        o.Debugs = d
    }
}

// SessionTimeout option
func SessionTimeout(t time.Duration) Option {
    return func(o *Options) {
        o.SessionTimeout = t
    }
}

// MaxPollInterval option
func MaxPollInterval(t time.Duration) Option {
    return func(o *Options) {
        o.MaxPollInterval = t
    }
}

// GroupInstanceID option
func GroupInstanceID(s string) Option {
    return func(o *Options) {
        o.GroupInstanceID = s
    }
}

// ConfigMapKey option
func ConfigMapKey(s string) Option {
    return func(o *Options) {
        o.ConfigMapKey = s
    }
}
