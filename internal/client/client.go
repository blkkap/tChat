package client

import (
	"fmt"	
	"os/signal"
	"log"
	"os"
	"bufio"
	"path/filepath"
	"github.com/gorilla/websocket"
	"encoding/json"
	"chat/internal/config"
)

var Cfg clientConfig
var in = bufio.NewReader(os.Stdin)

type Message struct{
        Message string `json:"message"`
	Username string `json:"username"`
}
type clientConfig struct{
	URL string `json:"URL"`
	USERNAME string `json:"USERNAME"`
}



func ensureConfigDir(){
	dir, err := os.UserConfigDir()
	if err != nil{
		log.Fatal(err)
	}
	path := filepath.Join(dir, "tchat")

	err = os.MkdirAll(filepath.Join(path,"client"), 0755)
	if err != nil{
		log.Fatal(err)
	}
}

func ensureClientConfig(){
	path := config.ClientConfigPath()
	if _, err := os.Stat(path); os.IsNotExist(err){
		fmt.Print("First Time Set Up: Creating client config.....")
		prompIfEmpty()
	}
}


func prompIfEmpty(){
	if Cfg.URL == "" || Cfg.USERNAME == ""{
		fmt.Print("\nEnter URL (e.g, : wss://URL.com/ws): ")
		var url string
		fmt.Scanln(&url)
		Cfg.URL = url

		fmt.Print("Enter USERNAME: ")
		var username string
		fmt.Scanln(&username)
		Cfg.USERNAME = username

		saveConfig()
	}
}

func saveConfig(){
	path := config.ClientConfigPath()
	data, err := json.MarshalIndent(Cfg,""," ")
	if err != nil{
		log.Fatal(err)
	}
	os.WriteFile(path,data,0644)
}

func getConfig(){

	file,err := os.Open(config.ClientConfigPath())
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


func Run(){
	ensureConfigDir()
	ensureClientConfig()
	getConfig()
	prompIfEmpty()
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
