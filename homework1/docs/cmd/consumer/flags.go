package main

import flag "github.com/spf13/pflag"

type flags struct {
	topic string
}

func init() {
	flag.StringVar(&cliFlags.topic, "topic", "pvz.events-log", "topic to consume from")

	flag.Parse()
}
