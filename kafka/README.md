# kafka

```golang
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/PengShaw/GoUtilsKit/kafka"
	"github.com/PengShaw/GoUtilsKit/logger"
)

func main() {
	brokers := []string{"192.168.1.2:9092"}
	topics := []string{"test"}
	group := "demo"
	k := kafka.NewKafkaConsumerGroup(brokers, topics, group, "range")
	k.SetOffset(sarama.OffsetOldest)
	k.SetConsumeClaimFunc(func(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
		for message := range claim.Messages() {
			logger.Infof("[topic:%s] [partiton:%d] [offset:%d] [value:%s] [time:%v]",
				message.Topic, message.Partition, message.Offset, string(message.Value), message.Timestamp)
			session.MarkMessage(message, "")
		}
		return nil
	})

	c, err := k.Connect()
	if err != nil {
		logger.Panic(err)
	}
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		logger.Warnln("terminating: via signal")
	}
	c()
}
```
