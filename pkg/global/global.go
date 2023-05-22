package global

var MQHost = "amqp://guest:guest@localhost:5672/"
var EtcdHost = "localhost:2379"
var ServerHost = "http://127.0.0.1:8080"

type Policy int

const (
	ROUND_ROBIN Policy = 0
	AFFINITY    Policy = 1
)
