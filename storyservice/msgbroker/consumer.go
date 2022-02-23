package msgbroker

import (
	"encoding/json"
	"errors"
	"log"
	"storyservice/adapters"
	"storyservice/msgbroker/listeners"

	"github.com/streadway/amqp"
)

const queueName = "story_worker"
const exchangeName = "user_exchange"
const routingKey = "story_feed"

var Consumers [1]func() = [1]func(){
	StoryFeedListener,
}

type EventData struct {
	EventType string
	Data      map[string]interface{}
}

func StoryFeedListener() {
	ch, err := adapters.GetRabbitmqConn().Channel()
	if err != nil {
		log.Println("Failed to open a channel")
		log.Fatal(err)
	}
	defer ch.Close()
	createExchange(ch)
	q := createQueue(ch)
	bindQueue(ch, q)

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
			e := EventData{}
			json.Unmarshal(d.Body, &e)
			log.Println(e)
			if listener, err := callListenerByEventType(&e); err != nil {
				log.Fatal(err)
			} else {
				go listener.Handle()
			}
		}
	}()

	log.Println("listening...")
	<-forever
}

func callListenerByEventType(eventData *EventData) (listeners.Listener, error) {
	var listener listeners.Listener
	switch eventData.EventType {
	case "follow":
		listener = &listeners.UserFollowEventData{Data: eventData.Data}
	case "unfollow":
		listener = &listeners.UserUnfollowEventData{Data: eventData.Data}
	default:
		return nil, errors.New("invalid event type")
	}
	return listener, nil
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

func createQueue(ch *amqp.Channel) *amqp.Queue {
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
	return &q
}

func bindQueue(ch *amqp.Channel, q *amqp.Queue) {
	err := ch.QueueBind(
		q.Name,       // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil)

	if err != nil {
		log.Fatal(err)
	}
}
