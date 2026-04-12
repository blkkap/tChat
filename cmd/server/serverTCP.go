package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "os"
    "encoding/json"
)

type config struct{
	PORT string `json:"PORT"`
}

var clients = make(map[net.Conn]bool)
var messages = make(chan string)
var Cfg config




func handleConnections(conn net.Conn){

	scanner := bufio.NewScanner(conn)
	clients[conn] = true
	defer func(){
		delete(clients,conn)
		conn.Close()
	}()
	
	for scanner.Scan(){
		msg := scanner.Text()
		formatted := fmt.Sprintf("%s", msg)

		messages <- formatted
	}
}


func getconfig(){
	file, err := os.Open("config.json")
	if err != nil{
		log.Fatal(err)
	}
	defer file.Close()
	
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Cfg)
	if err != nil{
		fmt.Println("Err decoding JSON:", err)
		return
	}
}


func broadcaster(){
	for{
		msg := <- messages
		for conn := range clients{
			fmt.Fprintln(conn,msg)
		}
	}
}


func main() {
    getconfig()
    

    ln, err := net.Listen("tcp", Cfg.PORT)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Listening on port", Cfg.PORT)
    go broadcaster()
    for{
    	conn, err := ln.Accept()
    	if err != nil {
        	log.Fatal(err)
     	}
	go handleConnections(conn)
    }

}


