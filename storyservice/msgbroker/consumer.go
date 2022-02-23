package msgbroker

import (
	"encoding/json"
	"log"
	"storyservice/adapters"
	"storyservice/msgbroker/listeners"

	"github.com/streadway/amqp"
)

const queueName = "story_worker"
const exchangeName = "user_exchange"
const routingKey = "story_feed"

var Consumers [1]func() = [1]func(){
	UserFollowListener,
}

func UserFollowListener() {
	ch, err := adapters.GetRabbitmqConn().Channel()
	if err != nil {
		log.Println("Failed to open a channel")
		log.Fatal(err)
	}
	defer ch.Close()
	createExchange(ch)
	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		true,      // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatal(err)
	}
	err = ch.QueueBind(
		q.Name,       // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Println("Received a message")
			m := listeners.UserFollowEventData{}
			json.Unmarshal(d.Body, &m)
			log.Println(m)
			go m.Handle()
		}
	}()

	log.Println("listening...")
	<-forever
}

func createExchange(ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatal(err)
	}
}
