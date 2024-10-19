package consumer

import (
	"github.com/IBM/sarama"
)

// Option is a configuration callback
type Option interface {
	Apply(*sarama.Config) error
}

type optionFn func(*sarama.Config) error

func (fn optionFn) Apply(c *sarama.Config) error {
	return fn(c)
}

func WithReturnErrorsEnabled(isEnabled bool) Option {
	return optionFn(func(c *sarama.Config) error {
		c.Consumer.Return.Errors = isEnabled
		return nil
	})
}
