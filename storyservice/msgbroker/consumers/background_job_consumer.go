package consumers

import (
	"encoding/json"
	"errors"
	"log"
	"storyservice/adapters"
	"storyservice/msgbroker/events"
	"storyservice/msgbroker/listeners"

	"github.com/streadway/amqp"
)

type BackgroundJobConsumer struct {
	queue   *amqp.Queue
	channel *amqp.Channel
}

func NewBackgroundJobConsumer() *BackgroundJobConsumer {
	ch, err := adapters.GetRabbitmqConn().Channel()
	if err != nil {
		log.Println("Failed to open a channel")
		log.Fatal(err)
	}

	return &BackgroundJobConsumer{channel: ch}
}

func (b *BackgroundJobConsumer) createExchange() {
	err := b.channel.ExchangeDeclare(
		events.BackgroundJobExchangeName, // name
		"direct",                         // type
		true,                             // durable
		false,                            // auto-deleted
		false,                            // internal
		false,                            // no-wait
		nil,                              // arguments
	)
	if err != nil {
		log.Fatal(err)
	}
}

func (b *BackgroundJobConsumer) createQueue() {
	q, err := b.channel.QueueDeclare(
		events.BackgroundJobQueueName, // name
		false,                         // durable
		false,                         // delete when unused
		true,                          // exclusive
		false,                         // no-wait
		nil,                           // arguments
	)
	if err != nil {
		log.Fatal(err)
	}
	b.queue = &q
}

func (b *BackgroundJobConsumer) bindQueue() {
	err := b.channel.QueueBind(
		b.queue.Name,                     // queue name
		events.BackgroundJobRoutingKey,   // routing key
		events.BackgroundJobExchangeName, // exchange
		false,
		nil)

	if err != nil {
		log.Fatal(err)
	}
}

func BackgroundJobHandler() {
	consumer := NewBackgroundJobConsumer()
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
			if listener, err := callBackgroundJobByEventType(&e); err != nil {
				log.Fatal(err)
			} else {
				go listener.Handle()
			}
		}
	}()

	log.Println("listening...")
	<-forever
}

func callBackgroundJobByEventType(eventData *EventData) (listeners.Listener, error) {
	var listener listeners.Listener
	switch eventData.EventType {
	case "chapter_read_log":
		listener = &listeners.ChapterReadLogListener{Data: eventData.Data}
	default:
		return nil, errors.New("invalid event type")
	}
	return listener, nil
}
