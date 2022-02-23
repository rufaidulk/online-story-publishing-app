package events

const BackgroundJobExchangeName = "background_job_exchange"
const BackgroundJobRoutingKey = "background_job"
const BackgroundJobQueueName = "background_job_worker"

type RabbitmqProducer interface {
	CreateExchange() RabbitmqProducer
	PublishMessage()
}

type EventData struct {
	EventType string
	Data      map[string]interface{}
}
