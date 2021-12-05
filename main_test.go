package main

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"testing"
)

func TestSandbox(t *testing.T) {
	var c = kafka.ConfigMap{}
	c.SetKey("test.mock.num.brokers", 3)
	var p, err = kafka.NewProducer(&c)
	if err != nil {
		t.Error(err)
	}
	p.Close()


}
