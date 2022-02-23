package events

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChapterReadLogEvent struct {
	EventChannel *amqp.Channel
	EventData    []byte
}

func NewChapterReadLogEvent(ch *amqp.Channel, eventType string, userUuid string, storyId primitive.ObjectID) *ChapterReadLogEvent {
	msg := make(map[string]interface{})
	msg["user_uuid"] = userUuid
	msg["story_id"] = storyId

	data := &EventData{
		EventType: eventType,
		Data:      msg,
	}
	byteData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	return &ChapterReadLogEvent{
		EventChannel: ch,
		EventData:    byteData,
	}
}

func (u *ChapterReadLogEvent) CreateExchange() RabbitmqProducer {
	err := u.EventChannel.ExchangeDeclare(
		BackgroundJobExchangeName, // name
		"direct",                  // type
		true,                      // durable
		false,                     // auto-deleted
		false,                     // internal
		false,                     // no-wait
		nil,                       // arguments
	)

	if err != nil {
		log.Println("Failed to declare an exchange")
		log.Fatal(err)
	}

	return u
}

func (u *ChapterReadLogEvent) PublishMessage() {
	err := u.EventChannel.Publish(
		BackgroundJobExchangeName, // exchange
		BackgroundJobRoutingKey,   // routing key
		false,                     // mandatory
		false,                     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        u.EventData,
		})

	if err != nil {
		log.Println("Failed to publish the message")
		log.Fatal(err)
	}
}
