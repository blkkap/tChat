package main

import (
	"fmt"
	"log"
	"os"
	"github.com/gorilla/websocket"
	"encoding/json"
)

var Cfg config


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


func main(){
	getConfig()
	var url = Cfg.URL
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil{
		log.Fatal("dial:", err)
	}
	defer conn.Close()
	/*for{
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Msg to send")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(conn,text + "\n")
		message,_ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print()
	}*/ 
	message := Cfg.USERNAME + ": Hello"
	err = conn.WriteMessage(websocket.TextMessage, []byte(message))
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
}
