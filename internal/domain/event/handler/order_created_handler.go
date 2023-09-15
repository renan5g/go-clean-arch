package handler

import (
	"encoding/json"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/renan5g/go-clean-arch/pkg/events"
	"github.com/renan5g/go-clean-arch/pkg/rabbitmq"
)

type OrderCreatedHandler struct {
	RabbitMQChannel *amqp.Channel
}

func NewOrderCreatedHandler(rabbitMQChannel *amqp.Channel) *OrderCreatedHandler {
	return &OrderCreatedHandler{RabbitMQChannel: rabbitMQChannel}
}

func (h *OrderCreatedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Order created: %v\n", event.GetPayload())
	jsonOutput, _ := json.Marshal(event.GetPayload())

	msgRabbitmq := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}

	rabbitmq.Publish(h.RabbitMQChannel, msgRabbitmq, "amq.direct")
}
