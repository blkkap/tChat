# About goChat
- Built using GO 
- Simple terminal based chat server

### Building from SRC 
- Building from source folder
  - git clone repo
  - run `go rum cmd/tchat/main.go` from root folder
  - follow prompt

### Config Files
- How to set up config files
  - Server side
    - Port is the only thing needed for the serverConfig it should be typed as follows: `:8080` 
    - Youll also need the url that the clients can find you on E.g. if using ngrok : `abc-def-123.com` (This will be used for clientConfig)
      - Also note if youre on http or https(refer to Client Side)
  - Client side
    - URL is needed from server host input url as follows: `wss://abc-def-123.com/ws` 
      - `if server side is https use wss:// if http use ws://`
    - Username selectiom





