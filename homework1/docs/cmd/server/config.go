package main

import "homework1/internal/infra/kafka"

type config struct {
	kafka kafka.Config
}

func newConfig() config {
	return config{
		kafka: kafka.Config{
			Brokers: []string{
				"localhost:9092",
			},
		},
	}
}
