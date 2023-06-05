package global

var MQHost = "amqp://guest:guest@localhost:5672/"
var EtcdHost = "127.0.0.1:2379"
var ServerHost = "http://127.0.0.1:8080"
var ServerlessHost = "http://127.0.0.1:8090"

var HostNameIpPrefix = "127.111.111"
var NameServerIp = "127.0.0.1"

type Policy int

const (
	ROUND_ROBIN Policy = 0
	AFFINITY    Policy = 1
)
