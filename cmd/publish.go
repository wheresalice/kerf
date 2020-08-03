/*
Copyright © 2020 wheresalice

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/cobra"
)

var broker string
var messages int
var topic string
var value string
var protocol string
var username string
var password string
var partition int

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish lots of Kafka messages",
	Long: `Publishes lots of messages to Kafka.`,
	Run: func(cmd *cobra.Command, args []string) {
		var c = kafka.ConfigMap{}
		c.SetKey("bootstrap.servers", broker)
		c.SetKey("linger.ms", 100)
		if protocol != "" {
			c.SetKey("sasl.mechanisms", "PLAIN")
			c.SetKey("security.protocol", protocol)
			c.SetKey("sasl.username", username)
			c.SetKey("sasl.password", password)
		}
		var p, err = kafka.NewProducer(&c)
		if err != nil {
			log.Printf("could not set up kafka producer: %s", err.Error())
			os.Exit(1)
		}
		done := make(chan bool)
		go func() {
			var msgCount int
			for e := range p.Events() {
				msg := e.(*kafka.Message)
				if msg.TopicPartition.Error != nil {
					log.Printf("delivery report error: %v", msg.TopicPartition.Error)
					os.Exit(1)
				}
				msgCount++
				if msgCount >= messages {
					done <- true
				}
			}
		}()

		defer p.Close()
		var start = time.Now()
		for j := 0; j < messages; j++ {
			p.ProduceChannel() <- &kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: int32(partition)}, Value: []byte(value)}
		}
		<-done
		elapsed := time.Since(start)

		log.Printf("[confluent-kafka-go producer] msg/s: %f", float64(messages) / elapsed.Seconds())
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)

	publishCmd.Flags().StringVarP(&broker, "broker", "b", "","bootstrap server")
	publishCmd.Flags().StringVarP(&topic, "topic", "t", "","kafka topic")
	publishCmd.Flags().StringVarP(&value, "value", "v", "","message value")
	publishCmd.Flags().IntVarP(&messages, "messages", "m", 1, "number of messages to generate")
	publishCmd.Flags().IntVarP(&partition, "partition", "p", 0, "topic partition")

	// Add some extra options to support Confluent Cloud
	// @todo surely we should be able to read these from a config file
	publishCmd.Flags().StringVarP(&protocol, "security.protocol", "s", "", "security.protocol")
	publishCmd.Flags().StringVarP(&username, "sasl.username", "u", "", "sasl.username")
	publishCmd.Flags().StringVarP(&password, "sasl.password", "k", "", "sasl.password")

	publishCmd.MarkFlagRequired("broker")
	publishCmd.MarkFlagRequired("topic")
	publishCmd.MarkFlagRequired("value")
}
