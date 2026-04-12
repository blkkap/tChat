package main

import(
	"fmt"
	"os"
	"log"
	"encoding/json"
	"net/http"
	"github.com/gorilla/websocket"
)


type config struct{
	PORT string `json:"PORT"`
}


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

	for{
		messageType, message, err := conn.ReadMessage()
		if err != nil{
			log.Println("read:", err)
			break
		}
		log.Printf("Recieved: %s", message)

		//Broaddcast as needed

		err = conn.WriteMessage(messageType, message)
		if err != nil{
			log.Println("Write:", err)
			break
		}
	}
}


func main(){
	getConfig()
	fmt.Println("Starting Server on: ", Cfg.PORT)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", websocketHandler)
	log.Fatal(http.ListenAndServe(Cfg.PORT, nil))
}
