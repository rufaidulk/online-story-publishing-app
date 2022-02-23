package consumers

import (
	"encoding/json"
	"errors"
	"log"
	"storyservice/adapters"
	"storyservice/msgbroker/listeners"

	"github.com/streadway/amqp"
)

type UserServiceConsumer struct {
	queue   *amqp.Queue
	channel *amqp.Channel
}

func NewUserServiceConsumer() *UserServiceConsumer {
	ch, err := adapters.GetRabbitmqConn().Channel()
	if err != nil {
		log.Println("Failed to open a channel")
		log.Fatal(err)
	}

	return &UserServiceConsumer{channel: ch}
}

func (u *UserServiceConsumer) createExchange() {
	err := u.channel.ExchangeDeclare(
		UserServiceExchangeName, // name
		"direct",                // type
		true,                    // durable
		false,                   // auto-deleted
		false,                   // internal
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		log.Fatal(err)
	}
}

func (u *UserServiceConsumer) createQueue() {
	q, err := u.channel.QueueDeclare(
		UserServiceQueueName, // name
		false,                // durable
		false,                // delete when unused
		true,                 // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		log.Fatal(err)
	}
	u.queue = &q
}

func (u *UserServiceConsumer) bindQueue() {
	err := u.channel.QueueBind(
		u.queue.Name,            // queue name
		UserServiceRoutingKey,   // routing key
		UserServiceExchangeName, // exchange
		false,
		nil)

	if err != nil {
		log.Fatal(err)
	}
}

func UserServiceListener() {
	consumer := NewUserServiceConsumer()
	defer consumer.channel.Close()
	consumer.createExchange()
	consumer.createQueue()
	consumer.bindQueue()

	msgs, err := consumer.channel.Consume(
		consumer.queue.Name, // queue
		"",                  // consumer
		true,                // auto ack
		false,               // exclusive
		false,               // no local
		false,               // no wait
		nil,                 // args
	)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Println("Received a message")
			e := EventData{}
			json.Unmarshal(d.Body, &e)
			log.Println(e)
			if listener, err := callUserServiceListenerByEventType(&e); err != nil {
				log.Fatal(err)
			} else {
				go listener.Handle()
			}
		}
	}()

	log.Println("listening...")
	<-forever
}

func callUserServiceListenerByEventType(eventData *EventData) (listeners.Listener, error) {
	var listener listeners.Listener
	switch eventData.EventType {
	case "follow":
		listener = &listeners.UserFollowEventListener{Data: eventData.Data}
	case "unfollow":
		listener = &listeners.UserUnfollowEventListener{Data: eventData.Data}
	default:
		return nil, errors.New("invalid event type")
	}
	return listener, nil
}
