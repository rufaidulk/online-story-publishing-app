package events

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type UserFollowEvent struct {
	EventChannel *amqp.Channel
	EventData    []byte
}

func NewUserFollowEvent(ch *amqp.Channel, eventType string, userUuid string, followerUuid string) *UserFollowEvent {
	msg := make(map[string]interface{})
	msg["user_uuid"] = userUuid
	msg["follower_uuid"] = followerUuid

	data := &EventData{
		EventType: eventType,
		Data:      msg,
	}
	byteData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	return &UserFollowEvent{
		EventChannel: ch,
		EventData:    byteData,
	}
}

func (u *UserFollowEvent) CreateExchange() RabbitmqProducer {
	err := u.EventChannel.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		log.Println("Failed to declare an exchange")
		log.Fatal(err)
	}

	return u
}

func (u *UserFollowEvent) PublishMessage() {
	err := u.EventChannel.Publish(
		exchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        u.EventData,
		})

	if err != nil {
		log.Println("Failed to publish the message")
		log.Fatal(err)
	}
}
