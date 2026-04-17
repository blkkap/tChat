package server

import(
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"context"
	"sync"
	"time"
	"log"
	"strconv"
	"encoding/json"
	"net/http"
	"path/filepath"
	"github.com/gorilla/websocket"
	"chat/internal/config"
)


type Serverconfig struct{
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
var Cfg Serverconfig


var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool {return true},

}

func ensureConfigDir(){
	
	dir, err := os.UserConfigDir()
	if err != nil{
		log.Fatal(err)
	}
	path := filepath.Join(dir, "tchat")
	
	err = os.MkdirAll(path, 0755)
	if err != nil{
		log.Fatal(err)
	}
	err = os.MkdirAll(filepath.Join(path, "server"), 0755)
	if err != nil{
		log.Fatal(err)
	}
}

func ensureServerConfig(){
	path := config.ServerConfigPath()

	if _, err := os.Stat(path); os.IsNotExist(err){
		fmt.Println("First Time Set Up: Creating server config....")
		promptIfEmpty()
	}
}

func getConfig(){

	file,err := os.Open(config.ServerConfigPath())
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

func promptIfEmpty(){
	if Cfg.PORT == ""{
		fmt.Print("Enter server port (e.g. :8080): ")

		var input string
		fmt.Scanln(&input)
		Cfg.PORT = input

		saveConfig()
	}
}


func saveConfig(){
	path := config.ServerConfigPath()

	data, err := json.MarshalIndent(Cfg, "", " ")
	if err != nil{
		log.Fatal(err)
	}
	os.WriteFile(path,data,0644)
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


func cleanDir(){
	path := config.ServerConfigPath()
	err := os.Remove(path)
	if err != nil{
		log.Fatal(err)
	}
}


func Run(){
	ensureConfigDir()
	ensureServerConfig()
	getConfig()
	promptIfEmpty()
	go broadcastMessages()
	//fmt.Println("Starting Server on: ", Cfg.PORT)


	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/ws", websocketHandler) 
	server := &http.Server{
		Addr: Cfg.PORT,
		Handler: mux,
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	
	go func(){
		log.Println("Starting Server on: ", Cfg.PORT)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed{
			log.Fatal("ListenAndServe: %v", err)
		}
	}()
	<-stop
	log.Println("self destruct......")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	if err := server.Shutdown(ctx); err != nil{
		log.Fatal("destruction failed:%+v", err)
	}
	log.Println("Stopped")
}
