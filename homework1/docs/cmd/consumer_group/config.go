package main

import "homework1/internal/infra/kafka"

type config struct {
	KafkaConfig kafka.Config
}

func newConfig(f flags) config {
	return config{
		KafkaConfig: kafka.Config{
			Brokers: []string{
				f.bootstrapServer,
			},
		},
	}
}
