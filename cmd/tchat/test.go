package main

import(
	"os"
	"fmt"
	"path/filepath"
)




func main(){
	fir, _ := os.UserConfigDir()

	path := filepath.Join(fir, "tchat")

	fmt.Print(path)
}
