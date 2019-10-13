package proxies

import (
	"fmt"
	"hi_bot/executor"
	"time"
)

type TestProxy struct{}

func (tp *TestProxy) Run(c chan *executor.Message, r chan []string) {
	var message *executor.Message
	var responseMessages []string

	time.Sleep(100 * time.Millisecond)
	fmt.Println("Proxy send message")
	message = executor.NewMessage("Test", "!Bot Test")
	c <- message
	responseMessages = <-r
	fmt.Println("Result", responseMessages)

	time.Sleep(400 * time.Millisecond)
	fmt.Println("Proxy send message")
	message = executor.NewMessage("Test", "!Bot OhOhoH")
	c <- message
	responseMessages = <-r
	fmt.Println("Result", responseMessages)

	time.Sleep(800 * time.Millisecond)
	fmt.Println("Proxy send message")
	message = executor.NewMessage("Test", "!Bot Yes")
	c <- message
	responseMessages = <-r
	fmt.Println("Result", responseMessages)
}
