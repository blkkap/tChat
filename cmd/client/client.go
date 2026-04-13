package main

import (
	
	"os/signal"
	"log"
	"os"
	"bufio"
	"github.com/gorilla/websocket"
	"encoding/json"
)

var Cfg config
var in = bufio.NewReader(os.Stdin)


type config struct{
	URL string `json:"URL"`
	USERNAME string `json:"USERNAME"`
}

func getConfig(){
	file,err := os.Open("config.json")
	if err != nil{
		log.Fatal(err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Cfg)
	if err != nil{
		log.Fatal(err)
	}
}

// Func that gets users input
func getInput(input chan string){
	results, err := in.ReadString('\n')
	if err != nil{
		log.Println(err)
		return
	}
	input <- results
}


func main(){
	getConfig()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	input := make(chan string, 1)
	go getInput(input)
	var url = Cfg.URL
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil{
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	done := make(chan struct{})
	go func(){
		defer close(done)
		for{
			_, message, err := conn.ReadMessage()
			if err != nil{
				log.Println("ReadMessage() err:", err)
				return
			}
			log.Printf("Recieved: %s", message)
		}
	}()
	for {
		select{
		case <- done:
			return
		case t := <-input:
			err := conn.WriteMessage(websocket.TextMessage, []byte(t))
			if err != nil{
				log.Println("write err:", err)
				return 
			}
			go getInput(input)
		case <-interrupt:
			log.Println("Caught interrupt signal")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

			if err != nil{
				log.Println("Write close err:", err)
				return
			}
			return
		}	

	}
	//message := Cfg.USERNAME + ": Hello"
	//err = conn.WriteMessage(websocket.TextMessage, []byte(message))
	/*
	if err != nil{
		log.Fatal(err)
	}

	for{
		_, message, err := conn.ReadMessage()
		if err != nil{
			log.Println("Read:", err)
			break
		}
		fmt.Printf("[%s]: %s\n",Cfg.USERNAME, message)
	}
	*/
}
