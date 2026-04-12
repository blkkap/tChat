package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "strings"
    "os"
    "encoding/json"
)

type config struct{
	PORT string `json:"PORT"` 
}

func main() {
    var Cfg config

    file, err := os.Open("../config.json")
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

    ln, err := net.Listen("tcp", Cfg.PORT)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Listening on port", Cfg.PORT)
    conn, err := ln.Accept()
    if err != nil {
        log.Fatal(err)
    }
    for {
        message, err :=  bufio.NewReader(conn).ReadString('\n')
        if err != nil {
            log.Fatal(err)
        }
        fmt.Print("Message Received:", string(message))
        newmessage := strings.ToUpper(message)
        conn.Write([]byte(newmessage + "\n"))
    }
}
