package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error){
	channel, err := conn.Channel()
	if(err != nil){
		return nil, amqp.Queue{}, err
	}

	durable := false
	autoDelate := false
	exclusive := false

	if(simpleQueueType == 0){
		durable = true
	}else{
		autoDelate = true
		exclusive = true
	}

	queue, err := channel.QueueDeclare(queueName, durable, autoDelate, exclusive, false, nil)
	if(err != nil){
		return nil, amqp.Queue{}, err
	}

	err = channel.QueueBind(queueName, key, exchange, false, nil)
	if(err != nil){
		return nil, amqp.Queue{}, err
	}

	return channel, queue, nil
}
