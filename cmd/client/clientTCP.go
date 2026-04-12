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
	go func(){
		
		reader := bufio.NewScanner(conn)
		for reader.Scan(){
			fmt.Println(reader.Text())
		}	
	}()
	message := bufio.NewScanner(os.Stdin)
	for message.Scan(){
		fmt.Println(message.Text())
	}
}
