package client

import (
	"log"
	"time"

	"github.com/aflog/assignment-messagebird/message"
	messagebird "github.com/messagebird/go-rest-api"
)

type mbClient struct {
	mc       *messagebird.Client
	throttle *time.Ticker
}

func NewMessageBird(key string) *mbClient {
	return &mbClient{
		mc:       messagebird.New(key),
		throttle: time.NewTicker(1 * time.Second),
	}
}

func (c *mbClient) Send(m message.Message) error {
	go func(msg message.Message) {
		<-c.throttle.C
		var params *messagebird.MessageParams
		if msg.UDH != "" {
			params = new(messagebird.MessageParams)
			params.Type = "binary"
			params.TypeDetails = map[string]interface{}{
				"udh": msg.UDH,
			}
		}
		res, err := c.mc.NewMessage(
			m.Originator,
			[]string{m.Recipient},
			m.Message,
			params)
		if err != nil {
			log.Printf("coul not send message %#v to MessageBird err: %s", msg, err.Error())
		}
		log.Printf("%#v", res)
		log.Printf("%#v", params)
		log.Printf("%#v", msg)
	}(m)
	return nil
}

func (c *mbClient) Close() {
	c.throttle.Stop()
}
