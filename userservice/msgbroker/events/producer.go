package events

// const queueName = "story_worker"
const exchangeName = "user_exchange"
const routingKey = "story_feed"

type RabbitmqProducer interface {
	CreateExchange() RabbitmqProducer
	PublishMessage()
}
