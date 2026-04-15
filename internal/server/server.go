package server

import(
	"fmt"
	"os"
	"sync"
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
	Username string `json:"username"`
	Message string `json:"message"`
	Sender *websocket.Conn
}

type Client struct{
	Conn *websocket.Conn
	Username string
}

var mu sync.Mutex
var clients = make(map[string]*Client)
var broadcast = make(chan Message)
var Cfg config


var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool {return true},

}



func getConfig(){

	file,err := os.Open("internal/config/serverConfig.json")
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
	//clients[clientID] = conn
	mu.Lock()
	clients[clientID] = &Client{
		Conn: conn,
	}
	mu.Unlock()
	defer delete(clients, clientID)

	log.Printf("Client connected: %s", clientID)


	for{
		/*
		_, message, err := conn.ReadMessage()
		if err != nil{
			log.Printf("read from client %s: %v:",clientID, err) //Log bug fixed
			break
		}
		log.Printf("Recieved from client %s: %s", clientID, message)
		*/	
		var msg Message 
		err := conn.ReadJSON(&msg)
		if err != nil{
			log.Printf("read from client %s: %v:", clientID, err)
			break
		}
		log.Printf("Received from [%s] : %s", msg.Username, msg.Message)
		//Broaddcast as needed

		broadcast <- Message{
			Username: msg.Username,
			Message: msg.Message,
			Sender: conn,
		}

		}
		closeError := conn.Close()
		if closeError != nil{
			log.Printf("Err closing connection for client %s: %v", clientID, closeError)
	}
}

func broadcastMessages(){

	for msg := range broadcast{
		formatted := fmt.Sprintf("%s : %s", msg.Username, msg.Message)
		mu.Lock()
		for id, client := range clients{
			if client.Conn == msg.Sender{
				continue
			}
			err := client.Conn.WriteMessage(websocket.TextMessage, []byte(formatted))
			if err != nil{
				log.Printf("err broadcasting to client %s: %v", id, err)
				client.Conn.Close()
				//mu.Lock()
				delete(clients,id)
				//mu.Unlock()
				//client.Close()
				}
			}
			mu.Unlock()
	}
}





func Run(){
	go broadcastMessages()
	getConfig()
	fmt.Println("Starting Server on: ", Cfg.PORT)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", websocketHandler) //wss Secure connection
	log.Fatal(http.ListenAndServe(Cfg.PORT, nil))
}
