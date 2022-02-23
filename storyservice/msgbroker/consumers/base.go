package consumers

const UserServiceQueueName = "user_service_worker"
const UserServiceExchangeName = "user_service_exchange"
const UserServiceRoutingKey = "user_service"

var Consumers [1]func() = [1]func(){
	UserServiceListener,
}

type EventData struct {
	EventType string
	Data      map[string]interface{}
}
