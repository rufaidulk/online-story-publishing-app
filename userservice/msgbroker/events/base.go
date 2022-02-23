package events

const ServiceExchangeName = "user_service_exchange"
const ServiceRoutingKey = "user_service"

type RabbitmqProducer interface {
	CreateExchange() RabbitmqProducer
	PublishMessage()
}

type EventData struct {
	EventType string
	Data      map[string]interface{}
}
