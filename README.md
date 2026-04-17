# About tChat
- Built using GO 
- Simple terminal based chat server


## Install tChat
### Mac / Linux
`curl -sL https://raw.githubusercontent.com/blkkap/tchat/main/install.sh | bash`

### Windows
- Download from GitHub Releases and run ./tchat.ext

## Building from SRC 
  ### Building from source folder
  - git clone repo
  - Install golang
  - Run `go rum cmd/tchat/main.go` from root folder
  - Follow prompt

## Config Files
  ### How to set up config files
  - Server side
    - Port is the only thing needed for the serverConfig it should be typed as follows: `:8080` 
    - Youll also need the url that the clients can find you on E.g. if using ngrok : `abc-def-123.com` (This will be used for clientConfig)
      - Note* if youre on http or https(refer to Client Side)
  - Client side
    - URL is needed from server host input URL as follows: `wss://abc-def-123.com/ws` 
      - `if server side is https use wss:// if http use ws://`
    - Username selectiom


## Clean up
- After each connection is closed (client and server) each config file is cleaned/deleted
- New config setup is required everytime running the program





