package pubsub

import(
	amqp "github.com/rabbitmq/amqp091-go"
	"encoding/json"
	"fmt"
)

type Acktype int

const(
	AcktypeAck = iota
	AcktypeNackRequeue
	AcktypeNackDiscard
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int,
	handler func(T)Acktype,
) error{
	channel, _, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if(err != nil){
		return err
	}

	delivery, err := channel.Consume(queueName, "", false, false, false, false, nil)
	if(err != nil){
		return err
	}

	go func() {
		for msg := range delivery {
			var data T
			err := json.Unmarshal(msg.Body, &data)
			if(err != nil){
				continue
			}
			ack := handler(data)
			switch(ack){
			case AcktypeAck:
				msg.Ack(false)
				fmt.Println("Accepted")
			case AcktypeNackRequeue:
				msg.Nack(false, true)
				fmt.Println("Rejected Requeue")
			case AcktypeNackDiscard:
				msg.Nack(false, false)
				fmt.Println("Rejected Discard")
			}
		}
		
	}()

	return nil
}