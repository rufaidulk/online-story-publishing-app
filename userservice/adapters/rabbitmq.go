package adapters

import (
	"log"
	"sync"

	"github.com/streadway/amqp"
)

var rabbitmqConn *amqp.Connection
var rabbitmqConnector sync.Once

func GetRabbitmqConn() *amqp.Connection {
	rabbitmqConnector.Do(func() {
		initRabbitmq()
	})

	return rabbitmqConn
}

func initRabbitmq() {
	var err error
	rabbitmqConn, err = amqp.Dial("amqp://root:root@message_broker:5672/")
	if err != nil {
		log.Fatal(err)
	}
}
