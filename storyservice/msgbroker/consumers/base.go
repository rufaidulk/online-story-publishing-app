package consumers

const UserServiceQueueName = "user_service_worker"
const UserServiceExchangeName = "user_service_exchange"
const UserServiceRoutingKey = "user_service"

var Consumers [2]func() = [2]func(){
	UserServiceListener,
	BackgroundJobHandler,
}

type EventData struct {
	EventType string
	Data      map[string]interface{}
}
