# About goChat
- Built using GO 
- Simple terminal based chat server



# Basic SetUP 

Server -> NGROK
CLient, Host of server, users to that server
Some sort of gui that will auto generate a ngrok server and fill in needed info and generate sign in key for users to join.
Some sort of config file for the server if self hosted (plug and play)
  - Server dies as soon as host ends. No saved server(or can we save that key for future connections? maybe re generate new keys every 24hr)


Out of Scope For Now:
- DB (Saved chats)
- Port to nvim package maybe, for now terminal only


# Things to Consider
- When host starts up server again if hosted via ngrok hostport will be newly assigned
  - If self hosted all stays the same
- SUB based? if it kicks off


# To Do LIST:
- Add config file  ✅
- add username to each user ✅ 
- extend code to hold multiple users (Server side done) ✅
- Lets user type out message which goes to server and everyone else except sender (SOMEWHAT)✅
- test with two users ✅
- Need to filter out sender own message ✅ 
- Clean config file ✅
- Set up user config files
- Add server auth
- Secure the server connection (sanatize)




## Last Things TODO
- Bundle app for easy installation
  - make new git repo for bundle(public)

## Least Important
- Set up command line arguments
- Figure out how to secure a log in for each server
- client side has echo response fix bug (NOT REALLY PRIORITY)
