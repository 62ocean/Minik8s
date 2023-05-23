package subscriber

import (
	"fmt"
)

type handler struct {
}

func (h handler) Handle(msg []byte) {
	fmt.Println("receive msg: " + string(msg))
}

func main() {
	s, _ := NewSubscriber("amqp://guest:guest@localhost:5672/")
	h := handler{}
	s.Subscribe("nodeQueue", Handler(h), nil)
}
