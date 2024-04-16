package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"github.com/PengShaw/GoUtilsKit/logger"
)

type KafkaConsumerGroup struct {
	brokers           []string
	topics            []string
	group             string
	assignor          sarama.BalanceStrategy
	offset            int64
	version           sarama.KafkaVersion
	channelBufferSize int
	ready             chan bool
	setupFunc         func(session sarama.ConsumerGroupSession) error
	cleanupFunc       func(sarama.ConsumerGroupSession) error
	consumeClaimFunc  func(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error
}

func (k *KafkaConsumerGroup) Brokers() []string {
	return k.brokers
}

func (k *KafkaConsumerGroup) Topics() []string {
	return k.topics
}

func (k *KafkaConsumerGroup) Group() string {
	return k.group
}

func (k *KafkaConsumerGroup) Assignor() sarama.BalanceStrategy {
	return k.assignor
}

func (k *KafkaConsumerGroup) Offset() int64 {
	return k.offset
}

func (k *KafkaConsumerGroup) Version() string {
	return k.version.String()
}

func (k *KafkaConsumerGroup) ChannelBufferSize() int {
	return k.channelBufferSize
}

func (k *KafkaConsumerGroup) SetOffset(offset int64) {
	k.offset = offset
}

func (k *KafkaConsumerGroup) SetVersion(version string) {
	k.version = getVersion(version)
}

func (k *KafkaConsumerGroup) SetChannelBufferSize(channelBufferSize int) {
	k.channelBufferSize = channelBufferSize
}

func (k *KafkaConsumerGroup) SetSetupFunc(setupFunc func(session sarama.ConsumerGroupSession) error) {
	k.setupFunc = setupFunc
}

func (k *KafkaConsumerGroup) SetCleanupFunc(cleanupFunc func(sarama.ConsumerGroupSession) error) {
	k.cleanupFunc = cleanupFunc
}
func (k *KafkaConsumerGroup) SetConsumeClaimFunc(consumeClaimFunc func(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error) {
	k.consumeClaimFunc = consumeClaimFunc
}

func NewKafkaConsumerGroup(brokers, topics []string, group string, assignor string) *KafkaConsumerGroup {
	return &KafkaConsumerGroup{
		brokers:           brokers,
		topics:            topics,
		group:             group,
		assignor:          getAssignor(assignor),
		offset:            sarama.OffsetNewest,
		version:           getVersion("3.7.0"),
		channelBufferSize: 1000,
		ready:             make(chan bool),
		setupFunc:         func(session sarama.ConsumerGroupSession) error { return nil },
		cleanupFunc:       func(session sarama.ConsumerGroupSession) error { return nil },
		consumeClaimFunc:  func(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error { return nil },
	}
}

func getAssignor(assignor string) sarama.BalanceStrategy {
	var a sarama.BalanceStrategy
	switch assignor {
	case "sticky":
		a = sarama.NewBalanceStrategySticky()
	case "roundrobin":
		a = sarama.NewBalanceStrategyRoundRobin()
	case "range":
		a = sarama.NewBalanceStrategyRange()
	default:
		logger.Panicf("consumer group partition assignor should be one of sticky, roundrobin or range, but got: %s", assignor)
	}
	return a
}

func getVersion(version string) sarama.KafkaVersion {
	v, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		logger.Panicf("Error parsing Kafka version: %v", err)
	}
	return v
}

func (k *KafkaConsumerGroup) Connect() (func(), error) {
	config := sarama.NewConfig()
	config.Version = k.version
	config.Consumer.Offsets.Initial = k.offset
	config.ChannelBufferSize = k.channelBufferSize

	client, err := sarama.NewConsumerGroup(k.brokers, k.group, config)
	if err != nil {
		return nil, fmt.Errorf("create consumer group client failed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, k.topics, k); err != nil {
				// when setup failed
				logger.Errorf("Error from consumer: %v", err)
				return
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				logger.Debugln("cancel connect: ", ctx.Err())
				return
			}
			k.ready <- true
		}
	}()
	<-k.ready
	logger.Infoln("kafka consumer group set up and running!...")

	return func() {
		logger.Info("kafka close")
		cancel()
		wg.Wait()
		close(k.ready)
		if err = client.Close(); err != nil {
			logger.Errorf("Error closing client: %v", err)
		}
	}, nil
}

func (k *KafkaConsumerGroup) Setup(session sarama.ConsumerGroupSession) error {
	logger.Debug("setup")
	return k.setupFunc(session)
}

func (k *KafkaConsumerGroup) Cleanup(session sarama.ConsumerGroupSession) error {
	logger.Debug("cleanup")
	return k.cleanupFunc(session)
}

func (k *KafkaConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	logger.Debug("consume claim")
	return k.consumeClaimFunc(session, claim)
}
