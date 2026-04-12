package main

import(
	"fmt"
	"os"
	"time"
	"log"
	"strconv"
	"encoding/json"
	"net/http"
	"github.com/gorilla/websocket"
)


type config struct{
	PORT string `json:"PORT"`
}

type Message struct{
	message []byte
	sender string
}


var clients = make(map[string]*websocket.Conn)
var broadcast = make(chan Message)
var Cfg config


var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool {return true},

}



func getConfig(){

	file, err := os.Open("config.json")
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

func homeHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome")

}

func websocketHandler(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w,r,nil)
	if err != nil{
		log.Println(err)
		return 
	}
	defer conn.Close()

	// Generate a unique id for clients
	clientID := strconv.Itoa(time.Now().Second())
	clients[clientID] = conn
	defer delete(clients, clientID)

	log.Printf("Client connected: %s", clientID)


	for{
		_, message, err := conn.ReadMessage()
		if err != nil{
			log.Println("read from client %s: %v:",clientID, err)
			break
		}
		log.Printf("Recieved from client %s: %s", clientID, message)

		//Broaddcast as needed

		broadcast <- Message{message, clientID}

		}
		closeError := conn.Close()
		if closeError != nil{
			log.Printf("Err closing connection for client %s: %v", clientID, closeError)
	}
}

func broadcastMessages(){
	for msg := range broadcast{
		for id, client := range clients{
			if id == msg.sender{
				continue
			}
			err := client.WriteMessage(websocket.TextMessage, msg.message)
			if err != nil{
				log.Printf("err broadcasting to client %s: %v", id, err)
				delete(clients,id)
				client.Close()
			}
		}
	}
}


func main(){
	go broadcastMessages()
	getConfig()
	fmt.Println("Starting Server on: ", Cfg.PORT)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", websocketHandler)
	log.Fatal(http.ListenAndServe(Cfg.PORT, nil))
}
