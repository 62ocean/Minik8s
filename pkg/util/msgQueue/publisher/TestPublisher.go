package publisher

import (
	"strconv"
	"time"
)

func main() {
	p, _ := NewPublisher("amqp://guest:guest@localhost:5672/")
	var i int
	for i = 0; ; i++ {
		p.Publish("nodeQueue", []byte("hello subscriber!"+strconv.Itoa(i)), "HELLO")
		time.Sleep(5 * time.Second)
	}
}
