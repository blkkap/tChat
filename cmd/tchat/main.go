package main

import(
	"fmt"
	"chat/internal/client"
	"chat/internal/server"
)


func main(){

	fmt.Println("1 : Start Client")
	fmt.Println("2 : Start Server")
	fmt.Println("\033[1 q")

	var choice string
	fmt.Scanln(&choice)

	switch choice{
		case "1":
			client.Run()
		case "2":
			server.Run()
		default:
			fmt.Println("Invalid choice")
	}
}
