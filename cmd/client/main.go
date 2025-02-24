package main

import (
	"fmt"
	"os"
	"os/signal"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)
	

func main() {
	connectionString := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectionString)
	if(err != nil){
		fmt.Println("dial error")
		return
	}

	defer connection.Close()

	

	name, err := gamelogic.ClientWelcome()
	if(err != nil){
		fmt.Println("name err")
		return
	}

	queueName := fmt.Sprintf("%s.%s", routing.PauseKey, name)
	_, _, err = pubsub.DeclareAndBind(connection, routing.ExchangePerilDirect, queueName, routing.PauseKey, 1)
	if(err != nil){
		fmt.Printf("declare and bind err %s\n", err)
		return
	}

	//err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{true})

	fmt.Println("Starting Peril client...")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
