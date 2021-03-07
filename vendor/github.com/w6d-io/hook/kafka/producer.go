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
    "encoding/json"
    "fmt"
    "time"

    kgo "github.com/confluentinc/confluent-kafka-go/kafka"
)

// Producer Creation of an emitter and send any type of value to its topic in kafka
func (k *Kafka) Producer(messageKey string, messageValue interface{}, opts ...Option) error {
    log := logger.WithName("Producer")

    options := NewOptions(opts...)

    producerCM := &kgo.ConfigMap{
        "bootstrap.servers": k.BootstrapServer,
    }
    if options.AuthKafka {
        _ = producerCM.SetKey("sasl.mechanisms", options.Mechanisms)
        _ = producerCM.SetKey("security.protocol", options.Protocol)
        _ = producerCM.SetKey("bootstrap.servers", k.BootstrapServer)
        _ = producerCM.SetKey("sasl.username", k.Username)
        _ = producerCM.SetKey("sasl.password", k.Password)
    }

    p, err := kgo.NewProducer(producerCM)
    if err != nil {
        return fmt.Errorf("failed to create producer: %s", err)
    }
    defer p.Close()
    go func() {
        for e := range p.Events() {
            switch ev := e.(type) {
            case *kgo.Message:
                if ev.TopicPartition.Error != nil {
                    log.Error(ev.TopicPartition.Error, "Failed to deliver",
                        "stacktrace", ev.TopicPartition)
                } else {
                    log.Info("Successfully produced record",
                        "topic", *ev.TopicPartition.Topic,
                        "partition", ev.TopicPartition.Partition,
                        "offset", ev.TopicPartition.Offset)
                }
            }
        }
    }()

    message, err := json.Marshal(&messageValue)
    if err != nil {
        log.Error(err, "marshal failed")
        return err
    }
    if err := p.Produce(&kgo.Message{
        TopicPartition: kgo.TopicPartition{Topic: &k.Topic, Partition: kgo.PartitionAny},
        Key:            []byte(messageKey),
        Value:          message,
        Timestamp:      time.Now(),
    }, nil); err != nil {
        log.Error(err, "produce failed")
        return err
    }
    p.Flush(int(options.WriteTimeout / time.Millisecond))
    return nil
}
