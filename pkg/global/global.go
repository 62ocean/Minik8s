package global

var MQHost = "amqp://guest:guest@localhost:5672/"
var EtcdHost = "localhost:2379"
var ServerHost = "http://192.168.1.6:8080"

var HostNameIpPrefix = "127.111.111"
var NameServerIp = "192.168.198.130"

type Policy int

const (
	ROUND_ROBIN Policy = 0
	AFFINITY    Policy = 1
)
