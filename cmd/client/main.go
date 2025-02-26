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

	state := gamelogic.NewGameState(name)

	for{
		input := gamelogic.GetInput()
		if len(input) == 0{
			continue
		}

		switch input[0]{
			case "spawn":
				err = state.CommandSpawn(input)
				if(err != nil){
					fmt.Println(err)
				}
			
			case "move":
				_, err = state.CommandMove(input)
				if(err != nil){
					fmt.Println(err)
				}else{
					fmt.Println("move successful")
				}
			case "status":
				state.CommandStatus()
			case "help":
				gamelogic.PrintClientHelp()
			case "spam":
				fmt.Println("not allowed")
			case "quit":
				gamelogic.PrintQuit()
				break
			default:
				fmt.Println("invalid command")
		}
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
