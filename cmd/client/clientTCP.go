package main

import (

	"fmt"
	"log"
	"net"
	"os"
	"bufio"
	"encoding/json"
)

type config struct{
	PORT string `json:"PORT"`
	LOCALHOST string `json:"LOCALHOST"`
}


func main(){
	var Cfg config
	
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

	conn,err := net.Dial("tcp", Cfg.LOCALHOST)
	if err != nil{
		log.Fatal(err)
	}
	for{
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(conn, text + "\n")
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: " + message)
	}
}
