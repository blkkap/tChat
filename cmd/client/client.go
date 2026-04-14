package main

import (
	"fmt"	
	"os/signal"
	"log"
	"os"
	"bufio"
	"github.com/gorilla/websocket"
	"encoding/json"
)

var Cfg config
var in = bufio.NewReader(os.Stdin)

type Message struct{
        Message string `json:"message"`
	Username string `json:"username"`
}
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
	for{
		//fmt.Print("\033[5 q") // [BAR CURSOR] might not be supported by all terminals
		fmt.Printf("%s : \033[1 q", Cfg.USERNAME) // [BLOCK]
		//fmt.Print("\033[3 q") // [UNDERLINE]
	
		//fmt.Printf("[%s] \r: ", Cfg.USERNAME)
		//time.Sleep(500 * time.Millisecond)
		results, err := in.ReadString('\n')
		if err != nil{
			log.Println(err)
			return
		}
		input <- results
	}
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
			fmt.Printf("\r\033[2K")
			fmt.Printf("%s", message)
			fmt.Printf("%s : \033[1 q", Cfg.USERNAME)
		}
	}()
	for {
		select{
		case <- done:
			return
		case t := <-input:
			msg := Message{
				Username: Cfg.USERNAME,
				Message: t,
			}
			data, err := json.Marshal(msg)
			if err != nil {
				log.Println("json err:", err)
				return
			}
			conn.WriteMessage(websocket.TextMessage, data)
			
			//go getInput(input)
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
}
