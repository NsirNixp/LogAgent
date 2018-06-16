package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

var client sarama.SyncProducer

func InitKafka(kafkaAddress string) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err = sarama.NewSyncProducer([]string{kafkaAddress}, config)
	if err != nil {
		logs.Error("Fail to create producer.Please check out the configuration [kafka::kafka_address]")
		err = fmt.Errorf("Fail to create producer.Please check out the configuration [kafka::kafka_address],[error=%v]", err)
		return
	}

	logs.Debug("InitKafka success,[kafkaAddress=%v]", kafkaAddress)
	fmt.Printf("InitKafka success,[kafkaAddress=%v]\n", kafkaAddress)
	return
}

func SendMessageToKafka(msg, topic string) (err error) {
	producerMessage := &sarama.ProducerMessage{
		Value: sarama.StringEncoder(msg),
		Topic: topic,
	}

	pid, offset, err := client.SendMessage(producerMessage)
	if err != nil {
		logs.Error("Fail to send ProducerMessage to kafka.Please check the kafka configuration or kafka server")
		err = fmt.Errorf("Fail to send ProducerMessage to kafka.Please check the kafka configuration or kafka server,[error=%v]", err)
		return
	}
	// 由于是一直运行的程序所以没有必要关闭client
	// defer client.Close()

	fmt.Printf("send to kafka success,===> pid:%v,===> offset:%v\n", pid, offset)
	return
}
