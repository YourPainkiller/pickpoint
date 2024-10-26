package main

import flag "github.com/spf13/pflag"

type flags struct {
	topic             string
	bootstrapServer   string
	consumerGroupName string
}

func init() {
	flag.StringVar(&cliFlags.topic, "topic", "pvz.events-log", "topic to consume from")
	flag.StringVar(&cliFlags.bootstrapServer, "bootstrap-server", "localhost:9092", "kafka broker host and port")
	flag.StringVar(&cliFlags.consumerGroupName, "cg-name", "route256-consumer-group", "consumer group name")

	flag.Parse()
}
