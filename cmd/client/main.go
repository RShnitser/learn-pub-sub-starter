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

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState)pubsub.Acktype{
	return func(ps routing.PlayingState)pubsub.Acktype{
		defer fmt.Print("> ")
		gs.HandlePause(ps)
		return pubsub.AcktypeAck
	}
}

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove)pubsub.Acktype {
	return func(move gamelogic.ArmyMove)pubsub.Acktype {
		defer fmt.Print("> ")
		outcome := gs.HandleMove(move)
		if outcome == gamelogic.MoveOutComeSafe || outcome == gamelogic.MoveOutcomeMakeWar{
			return pubsub.AcktypeAck
		}else{
			return pubsub.AcktypeNackDiscard
		}
	}
}
	

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

	armyName := fmt.Sprintf("%s.%s", "army_moves", name)
	armyChannel, _, err := pubsub.DeclareAndBind(connection, routing.ExchangePerilTopic, armyName, "army_moves.*", 1)
	if(err != nil){
		fmt.Printf("declare and bind err %s\n", err)
		return
	}

	fmt.Println("Starting Peril client...")

	state := gamelogic.NewGameState(name)

	
	err = pubsub.SubscribeJSON(connection, routing.ExchangePerilDirect, queueName, routing.PauseKey, 1, handlerPause(state))
	if(err != nil){
		fmt.Printf("subscribe error", err)
		return
	}

	err = pubsub.SubscribeJSON(connection, routing.ExchangePerilTopic, armyName, "army_moves.*", 1, handlerMove(state))
	if(err != nil){
		fmt.Printf("subscribe error", err)
		return
	}

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
				move, err := state.CommandMove(input)
				if(err != nil){
					fmt.Println(err)
				}else{
					err = pubsub.PublishJSON(armyChannel, routing.ExchangePerilTopic, armyName, move)
					if(err != nil){
						fmt.Println(err)
					}else{	
							fmt.Println("move successful")
					}
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
