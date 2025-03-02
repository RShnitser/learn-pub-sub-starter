package main

import (
	"fmt"
	"os"
	"os/signal"
    amqp "github.com/rabbitmq/amqp091-go"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
)


func main() {
	connectionString := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectionString)
	if(err != nil){
		return
	}

	defer connection.Close()

	channel, err := connection.Channel()
	if(err != nil){
		return
	}

	_, _, err = pubsub.DeclareAndBind(connection, "peril_topic", "game_logs", "game_logs.*", 0)
	if(err != nil){
		fmt.Println("bind error")
		return
	}

	fmt.Println("Starting Peril server...")

	gamelogic.PrintServerHelp()

	for{
		input := gamelogic.GetInput()
		if len(input) == 0{
			continue
		}

		if input[0] == "pause"{
			err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{true})
		}else if input[0] == "resume"{
			err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{false})
		}else if(input[0] == "quit"){
			break
		}else{
			fmt.Println("invalid command")
		}
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("Shutting Down")
}
